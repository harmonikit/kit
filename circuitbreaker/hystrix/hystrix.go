// Package hystrix provides a Hystrix-style circuit breaker implementing
// harmoni/circuitbreaker.
//
// The circuit opens after a configurable number of consecutive failures
// and transitions to half-open after a sleep window.
package hystrix

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/harmonikit/harmoni/circuitbreaker"
	"github.com/harmonikit/harmoni/endpoint"
)

// Compile-time interface check.
var _ circuitbreaker.CircuitBreaker[int, int] = (*CircuitBreaker[int, int])(nil)

// Config holds the circuit breaker configuration.
type Config struct {
	// MaxConsecutiveFailures is the number of failures before opening.
	MaxConsecutiveFailures int
	// SleepWindow is how long the circuit stays open before going half-open.
	SleepWindow time.Duration
}

// DefaultConfig returns a reasonable default configuration.
func DefaultConfig() Config {
	return Config{
		MaxConsecutiveFailures: 5,
		SleepWindow:            30 * time.Second,
	}
}

// CircuitBreaker is a Hystrix-style circuit breaker.
type CircuitBreaker[Req, Resp any] struct {
	config   Config
	mu       sync.Mutex
	failures int
	state    circuitbreaker.State
	openedAt time.Time
}

// New returns a new Hystrix circuit breaker.
func New[Req, Resp any](cfg Config) *CircuitBreaker[Req, Resp] {
	if cfg.MaxConsecutiveFailures <= 0 {
		cfg.MaxConsecutiveFailures = 5
	}
	if cfg.SleepWindow <= 0 {
		cfg.SleepWindow = 30 * time.Second
	}
	return &CircuitBreaker[Req, Resp]{
		config:   cfg,
		mu:       sync.Mutex{},
		failures: 0,
		state:    circuitbreaker.StateClosed,
		openedAt: time.Time{},
	}
}

// Execute runs the endpoint through the circuit breaker.
func (cb *CircuitBreaker[Req, Resp]) Execute(ctx context.Context, req Req, ep endpoint.Endpoint[Req, Resp]) (Resp, error) {
	if !cb.allow() {
		var zero Resp
		return zero, fmt.Errorf("circuit breaker: %s", cb.state)
	}

	resp, err := ep(ctx, req)

	cb.record(err)
	return resp, err
}

// State returns the current state.
func (cb *CircuitBreaker[Req, Resp]) State() circuitbreaker.State {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.transition()
	return cb.state
}

func (cb *CircuitBreaker[Req, Resp]) allow() bool {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.transition()
	return cb.state != circuitbreaker.StateOpen
}

func (cb *CircuitBreaker[Req, Resp]) record(err error) {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	if err != nil {
		cb.failures++
	} else {
		cb.failures = 0
		if cb.state == circuitbreaker.StateHalfOpen {
			cb.state = circuitbreaker.StateClosed
		}
	}

	cb.transition()
}

func (cb *CircuitBreaker[Req, Resp]) transition() {
	switch cb.state {
	case circuitbreaker.StateClosed:
		if cb.failures >= cb.config.MaxConsecutiveFailures {
			cb.state = circuitbreaker.StateOpen
			cb.openedAt = time.Now()
		}
	case circuitbreaker.StateOpen:
		if time.Since(cb.openedAt) >= cb.config.SleepWindow {
			cb.state = circuitbreaker.StateHalfOpen
		}
	}
}
