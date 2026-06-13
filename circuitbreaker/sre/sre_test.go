package sre_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/harmonikit/harmoni/circuitbreaker"
	"github.com/harmonikit/harmoni/endpoint"
	cbsre "github.com/harmonikit/kit/circuitbreaker/sre"
)

func TestCircuitBreaker_Closed(t *testing.T) {
	cb := cbsre.New[int, int](cbsre.DefaultConfig())

	if cb.State() != circuitbreaker.StateClosed {
		t.Fatalf("expected closed, got %v", cb.State())
	}
}

func TestCircuitBreaker_Execute_Success(t *testing.T) {
	cb := cbsre.New[int, int](cbsre.DefaultConfig())

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
}

func TestCircuitBreaker_OpensOnHighErrorRate(t *testing.T) {
	cfg := cbsre.Config{
		ErrorThreshold: 0.5,
		WindowSize:     10,
		SleepWindow:    time.Hour,
	}
	cb := cbsre.New[int, int](cfg)

	failEP := endpoint.Endpoint[int, int](func(ctx context.Context, req int) (int, error) {
		return 0, errors.New("fail")
	})

	// Fill the window with failures.
	for range 10 {
		cb.Execute(context.Background(), 1, failEP)
	}

	if cb.State() != circuitbreaker.StateOpen {
		t.Fatalf("expected open after 100%% failures, got %v", cb.State())
	}
}

func TestCircuitBreaker_StaysClosedOnLowErrorRate(t *testing.T) {
	cfg := cbsre.Config{
		ErrorThreshold: 0.5,
		WindowSize:     10,
		SleepWindow:    time.Hour,
	}
	cb := cbsre.New[int, int](cfg)

	okEP := endpoint.Endpoint[int, int](func(ctx context.Context, req int) (int, error) {
		return req, nil
	})
	failEP := endpoint.Endpoint[int, int](func(ctx context.Context, req int) (int, error) {
		return 0, errors.New("fail")
	})

	// 70% success, 30% failure — below 50% threshold.
	for range 7 {
		cb.Execute(context.Background(), 1, okEP)
	}
	for range 3 {
		cb.Execute(context.Background(), 1, failEP)
	}

	if cb.State() != circuitbreaker.StateClosed {
		t.Fatalf("expected closed at 30%% error rate, got %v", cb.State())
	}
}

func TestCircuitBreaker_HalfOpenToClosed(t *testing.T) {
	cfg := cbsre.Config{
		ErrorThreshold: 0.5,
		WindowSize:     10,
		SleepWindow:    1 * time.Millisecond,
	}
	cb := cbsre.New[int, int](cfg)

	failEP := endpoint.Endpoint[int, int](func(ctx context.Context, req int) (int, error) {
		return 0, errors.New("fail")
	})
	for range 10 {
		cb.Execute(context.Background(), 1, failEP)
	}

	time.Sleep(10 * time.Millisecond)

	okEP := endpoint.Endpoint[int, int](func(ctx context.Context, req int) (int, error) {
		return req, nil
	})
	_, err := cb.Execute(context.Background(), 1, okEP)
	if err != nil {
		t.Fatalf("expected half-open to allow: %v", err)
	}
	if cb.State() != circuitbreaker.StateClosed {
		t.Fatalf("expected closed after half-open success, got %v", cb.State())
	}
}
