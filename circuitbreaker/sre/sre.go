// Package sre provides a Google SRE-style adaptive circuit breaker.
//
// It tracks the error ratio over a sliding window and opens the circuit
// when the ratio exceeds a configurable threshold.
package sre

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

// Config holds the SRE circuit breaker configuration.
type Config struct {
	ErrorThreshold float64
	WindowSize     int
	SleepWindow    time.Duration
}

// DefaultConfig returns a reasonable default configuration.
func DefaultConfig() Config {
	return Config{
		ErrorThreshold: 0.5,
		WindowSize:     100,
		SleepWindow:    30 * time.Second,
	}
}

// CircuitBreaker is an SRE-style adaptive circuit breaker.
type CircuitBreaker[Req, Resp any] struct {
	config   Config
	mu       sync.Mutex
	results  []bool
	idx      int
	state    circuitbreaker.State
	openedAt time.Time
}

// New returns a new SRE circuit breaker.
func New[Req, Resp any](cfg Config) *CircuitBreaker[Req, Resp] {
	if cfg.ErrorThreshold <= 0 {
		cfg.ErrorThreshold = 0.5
	}
	if cfg.WindowSize <= 0 {
		cfg.WindowSize = 100
	}
	if cfg.SleepWindow <= 0 {
		cfg.SleepWindow = 30 * time.Second
	}
	return &CircuitBreaker[Req, Resp]{
		config:   cfg,
		mu:       sync.Mutex{},
		results:  make([]bool, 0, cfg.WindowSize),
		idx:      0,
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

	cb.record(err == nil)
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

func (cb *CircuitBreaker[Req, Resp]) record(success bool) {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	if len(cb.results) < cb.config.WindowSize {
		cb.results = append(cb.results, success)
	} else {
		cb.results[cb.idx] = success
		cb.idx = (cb.idx + 1) % cb.config.WindowSize
	}

	cb.transition()
}

func (cb *CircuitBreaker[Req, Resp]) transition() {
	switch cb.state {
	case circuitbreaker.StateClosed:
		if cb.errorRatio() > cb.config.ErrorThreshold {
			cb.state = circuitbreaker.StateOpen
			cb.openedAt = time.Now()
		}
	case circuitbreaker.StateOpen:
		if time.Since(cb.openedAt) >= cb.config.SleepWindow {
			cb.state = circuitbreaker.StateHalfOpen
		}
	}
}

func (cb *CircuitBreaker[Req, Resp]) errorRatio() float64 {
	if len(cb.results) == 0 {
		return 0
	}
	var failures int
	for _, r := range cb.results {
		if !r {
			failures++
		}
	}
	return float64(failures) / float64(len(cb.results))
}
