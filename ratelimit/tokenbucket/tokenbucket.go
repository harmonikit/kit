// Package tokenbucket provides a production-grade token bucket rate limiter.
//
// It builds on harmoni/ratelimit.Limiter with additional features like
// dynamic rate adjustment and rate limit middleware for endpoints.
package tokenbucket

import (
	"context"
	"sync"
	"time"

	"github.com/harmonikit/harmoni/endpoint"
)

// Limiter is a production-grade token bucket rate limiter.
// It implements ratelimit.Limiter and adds dynamic rate adjustment.
type Limiter struct {
	mu       sync.Mutex
	rate     float64 // tokens per second
	burst    float64 // max tokens
	tokens   float64
	lastTime time.Time
}

// New creates a Limiter with the given rate and burst.
func New(rate float64, burst int) *Limiter {
	burstf := float64(burst)
	return &Limiter{
		rate:     rate,
		burst:    burstf,
		tokens:   burstf,
		lastTime: time.Now(),
	}
}

// Allow reports whether a request is allowed.
func (l *Limiter) Allow() bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(l.lastTime).Seconds()
	l.lastTime = now

	l.tokens += elapsed * l.rate
	if l.tokens > l.burst {
		l.tokens = l.burst
	}

	if l.tokens >= 1 {
		l.tokens--
		return true
	}
	return false
}

// SetRate dynamically adjusts the rate (tokens per second).
func (l *Limiter) SetRate(rate float64) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.rate = rate
}

// Middleware returns an endpoint middleware that rate-limits requests.
// When the limiter denies a request, the endpoint is not called and an
// error is returned.
func Middleware[Req, Resp any](limiter *Limiter) endpoint.Middleware[Req, Resp] {
	return func(next endpoint.Endpoint[Req, Resp]) endpoint.Endpoint[Req, Resp] {
		return func(ctx context.Context, req Req) (Resp, error) {
			if !limiter.Allow() {
				var zero Resp
				return zero, RateLimitError{}
			}
			return next(ctx, req)
		}
	}
}

// RateLimitError is returned when the rate limit is exceeded.
type RateLimitError struct{}

func (e RateLimitError) Error() string {
	return "rate limit exceeded"
}
