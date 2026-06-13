package hystrix_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/harmonikit/harmoni/circuitbreaker"
	"github.com/harmonikit/harmoni/endpoint"
	cbhystrix "github.com/harmonikit/kit/circuitbreaker/hystrix"
)

func TestCircuitBreaker_Closed(t *testing.T) {
	cb := cbhystrix.New[int, int](cbhystrix.DefaultConfig())

	if cb.State() != circuitbreaker.StateClosed {
		t.Fatalf("expected closed, got %v", cb.State())
	}
}

func TestCircuitBreaker_Execute_Success(t *testing.T) {
	cb := cbhystrix.New[int, int](cbhystrix.DefaultConfig())

	ep := endpoint.Endpoint[int, int](func(ctx context.Context, req int) (int, error) {
		return req * 2, nil
	})

	resp, err := cb.Execute(context.Background(), 21, ep)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp != 42 {
		t.Fatalf("got %d, want 42", resp)
	}
	if cb.State() != circuitbreaker.StateClosed {
		t.Fatalf("expected closed after success, got %v", cb.State())
	}
}

func TestCircuitBreaker_OpensAfterFailures(t *testing.T) {
	cfg := cbhystrix.Config{
		MaxConsecutiveFailures: 3,
		SleepWindow:            time.Hour,
	}
	cb := cbhystrix.New[int, int](cfg)

	ep := endpoint.Endpoint[int, int](func(ctx context.Context, req int) (int, error) {
		return 0, errors.New("fail")
	})

	for range 3 {
		cb.Execute(context.Background(), 1, ep)
	}

	if cb.State() != circuitbreaker.StateOpen {
		t.Fatalf("expected open after failures, got %v", cb.State())
	}
}

func TestCircuitBreaker_OpenBlocks(t *testing.T) {
	cfg := cbhystrix.Config{
		MaxConsecutiveFailures: 1,
		SleepWindow:            time.Hour,
	}
	cb := cbhystrix.New[int, int](cfg)

	failEP := endpoint.Endpoint[int, int](func(ctx context.Context, req int) (int, error) {
		return 0, errors.New("fail")
	})
	cb.Execute(context.Background(), 1, failEP)

	okEP := endpoint.Endpoint[int, int](func(ctx context.Context, req int) (int, error) {
		return req, nil
	})
	_, err := cb.Execute(context.Background(), 1, okEP)
	if err == nil {
		t.Fatal("expected error when circuit is open")
	}
}

func TestCircuitBreaker_HalfOpenToClosed(t *testing.T) {
	cfg := cbhystrix.Config{
		MaxConsecutiveFailures: 1,
		SleepWindow:            1 * time.Millisecond,
	}
	cb := cbhystrix.New[int, int](cfg)

	failEP := endpoint.Endpoint[int, int](func(ctx context.Context, req int) (int, error) {
		return 0, errors.New("fail")
	})
	cb.Execute(context.Background(), 1, failEP)

	// Wait for sleep window.
	time.Sleep(10 * time.Millisecond)

	okEP := endpoint.Endpoint[int, int](func(ctx context.Context, req int) (int, error) {
		return req, nil
	})
	_, err := cb.Execute(context.Background(), 1, okEP)
	if err != nil {
		t.Fatalf("expected half-open to allow request: %v", err)
	}
	if cb.State() != circuitbreaker.StateClosed {
		t.Fatalf("expected closed after success in half-open, got %v", cb.State())
	}
}
