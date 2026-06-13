package legacy_test

import (
	"context"
	"errors"
	"testing"

	"github.com/harmonikit/harmoni/endpoint"
	"github.com/harmonikit/kit/legacy"
)

func TestAdapt(t *testing.T) {
	oldEP := legacy.Endpoint(func(ctx context.Context, request any) (any, error) {
		return request.(int) * 2, nil
	})

	newEP := legacy.Adapt(oldEP)
	resp, err := newEP(context.Background(), 21)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.(int) != 42 {
		t.Fatalf("got %d, want 42", resp)
	}
}

func TestAdapt_Error(t *testing.T) {
	wantErr := errors.New("legacy error")
	oldEP := legacy.Endpoint(func(ctx context.Context, request any) (any, error) {
		return nil, wantErr
	})

	newEP := legacy.Adapt(oldEP)
	_, err := newEP(context.Background(), "test")
	if !errors.Is(err, wantErr) {
		t.Fatalf("got %v, want %v", err, wantErr)
	}
}

func TestAdaptTyped(t *testing.T) {
	oldEP := legacy.Endpoint(func(ctx context.Context, request any) (any, error) {
		return request.(int) * 3, nil
	})

	typedEP := legacy.AdaptTyped[int, int](oldEP)
	resp, err := typedEP(context.Background(), 7)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp != 21 {
		t.Fatalf("got %d, want 21", resp)
	}
}

func TestAdaptTyped_WrongType(t *testing.T) {
	oldEP := legacy.Endpoint(func(ctx context.Context, request any) (any, error) {
		return "string-response", nil
	})

	typedEP := legacy.AdaptTyped[int, int](oldEP)
	_, err := typedEP(context.Background(), 1)
	if err == nil {
		t.Fatal("expected type assertion error")
	}
}

func TestLegacy_Interface(t *testing.T) {
	var _ legacy.Endpoint = func(ctx context.Context, request any) (any, error) {
		return request, nil
	}

	// Verify Adapt returns a harmoni endpoint.
	var _ endpoint.Endpoint[any, any] = legacy.Adapt(nil)
}
