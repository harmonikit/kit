// Package authz provides authorization middleware on top of harmoni/auth.
//
// After harmoni/auth authenticates a request and embeds claims in the context,
// authz checks whether the request is authorized by evaluating policies.
//
// Example:
//
//	ep = endpoint.Chain(
//	    auth.Middleware[Req, Resp](jwtAuth),
//	    authz.Middleware[Req, Resp](authz.AllowAll[Req]()),
//	)(ep)
package authz

import (
	"context"
	"fmt"

	"github.com/harmonikit/harmoni/endpoint"
)

// Policy evaluates whether a request is authorized.
type Policy[Req any] interface {
	// Allowed returns nil if the request is authorized, or an error if not.
	Allowed(ctx context.Context, req Req) error
}

// Middleware returns an endpoint middleware that enforces a policy.
// It should be placed after authentication middleware in the chain.
func Middleware[Req, Resp any](p Policy[Req]) endpoint.Middleware[Req, Resp] {
	return func(next endpoint.Endpoint[Req, Resp]) endpoint.Endpoint[Req, Resp] {
		return func(ctx context.Context, req Req) (Resp, error) {
			if err := p.Allowed(ctx, req); err != nil {
				var zero Resp
				return zero, fmt.Errorf("authz: %w", err)
			}
			return next(ctx, req)
		}
	}
}

// AllowAll returns a policy that allows every request. Useful as a default
// or for testing.
type AllowAll[Req any] struct{}

func (AllowAll[Req]) Allowed(_ context.Context, _ Req) error { return nil }

// DenyAll returns a policy that denies every request.
type DenyAll[Req any] struct{}

func (DenyAll[Req]) Allowed(_ context.Context, _ Req) error {
	return fmt.Errorf("access denied")
}

// ErrUnauthorized is returned when authorization fails.
var ErrUnauthorized = fmt.Errorf("unauthorized")
