package authz_test

import (
	"context"
	"errors"
	"testing"

	"github.com/harmonikit/harmoni/endpoint"
	"github.com/harmonikit/kit/auth/authz"
)

func TestAllowAll(t *testing.T) {
	ep := endpoint.Endpoint[int, int](func(ctx context.Context, req int) (int, error) {
		return req * 2, nil
	})

	p := authz.AllowAll[int]{}
	wrapped := authz.Middleware[int, int](p)(ep)

	resp, err := wrapped(context.Background(), 21)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp != 42 {
		t.Fatalf("got %d, want 42", resp)
	}
}

func TestDenyAll(t *testing.T) {
	ep := endpoint.Endpoint[int, int](func(ctx context.Context, req int) (int, error) {
		return req, nil
	})

	p := authz.DenyAll[int]{}
	wrapped := authz.Middleware[int, int](p)(ep)

	_, err := wrapped(context.Background(), 1)
	if err == nil {
		t.Fatal("expected authorization error")
	}
}

func TestCustomPolicy(t *testing.T) {
	adminOnly := &rolePolicy{role: "admin"}
	ep := endpoint.Endpoint[string, string](func(ctx context.Context, req string) (string, error) {
		return "ok", nil
	})

	wrapped := authz.Middleware[string, string](adminOnly)(ep)

	// Request without role context should fail.
	_, err := wrapped(context.Background(), "request")
	if err == nil {
		t.Fatal("expected authorization error without role")
	}

	// Request with admin role should succeed.
	ctx := context.WithValue(context.Background(), "role", "admin")
	resp, err := wrapped(ctx, "request")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp != "ok" {
		t.Fatalf("got %q, want ok", resp)
	}
}

func TestMiddleware_Chains(t *testing.T) {
	// Verify that middleware composes correctly in a chain.
	p := authz.AllowAll[int]{}
	ep := endpoint.Endpoint[int, int](func(ctx context.Context, req int) (int, error) {
		return req, nil
	})

	wrapped := endpoint.Chain(authz.Middleware[int, int](p))(ep)
	_, err := wrapped(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error in chain: %v", err)
	}
}

// rolePolicy allows requests only when the context has a matching role.
type rolePolicy struct {
	role string
}

func (p *rolePolicy) Allowed(ctx context.Context, _ string) error {
	role := ctx.Value("role")
	if role == nil || role.(string) != p.role {
		return errors.New("access denied")
	}
	return nil
}
