package tokenbucket_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/harmonikit/harmoni/endpoint"
	tb "github.com/harmonikit/kit/ratelimit/tokenbucket"
)

func TestLimiter_Allow(t *testing.T) {
	lim := tb.New(1000, 3)

	for range 3 {
		if !lim.Allow() {
			t.Fatal("expected allow within burst")
		}
	}
	// Burst exhausted — should deny.
	if lim.Allow() {
		t.Fatal("expected deny after burst")
	}
}

func TestLimiter_Refill(t *testing.T) {
	lim := tb.New(10000, 1)

	if !lim.Allow() {
		t.Fatal("expected initial allow")
	}
	time.Sleep(1 * time.Millisecond)
	if !lim.Allow() {
		t.Fatal("expected allow after refill")
	}
}

func TestLimiter_SetRate(t *testing.T) {
	lim := tb.New(0, 2)

	// Consume both tokens.
	lim.Allow()
	lim.Allow()

	// Should be denied at zero rate.
	if lim.Allow() {
		t.Fatal("expected deny at zero rate")
	}

	// Increase rate.
	lim.SetRate(10000)
	time.Sleep(1 * time.Millisecond)

	if !lim.Allow() {
		t.Fatal("expected allow after rate increase")
	}
}

func TestMiddleware_Allows(t *testing.T) {
	lim := tb.New(10000, 10)
	mw := tb.Middleware[int, int](lim)

	ep := endpoint.Endpoint[int, int](func(ctx context.Context, req int) (int, error) {
		return req * 2, nil
	})

	wrapped := mw(ep)
	resp, err := wrapped(context.Background(), 21)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp != 42 {
		t.Fatalf("got %d, want 42", resp)
	}
}

func TestMiddleware_Denies(t *testing.T) {
	lim := tb.New(0, 0) // zero burst, zero rate
	mw := tb.Middleware[int, int](lim)

	ep := endpoint.Endpoint[int, int](func(ctx context.Context, req int) (int, error) {
		return req, nil
	})

	wrapped := mw(ep)
	_, err := wrapped(context.Background(), 1)
	if err == nil {
		t.Fatal("expected rate limit error")
	}
	var rle tb.RateLimitError
	if !errors.As(err, &rle) {
		t.Fatalf("got %T, want RateLimitError", err)
	}
}
