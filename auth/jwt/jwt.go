// Package jwt provides JWT-based authentication implementing harmoni/auth.Auth.
//
// This is a stub implementation. A production version depends on
// github.com/golang-jwt/jwt/v5 for token parsing and validation.
package jwt

import (
	"context"
	"errors"
	"fmt"

	"github.com/harmonikit/harmoni/auth"
)

// Auth implements auth.Auth for JWT-based authentication.
// It extracts and validates a JWT token from the request.
type Auth[Req any] struct {
	// ExtractToken extracts the JWT token string from the request.
	ExtractToken func(ctx context.Context, req Req) (string, error)
	// ValidateToken validates the token and returns claims.
	ValidateToken func(ctx context.Context, tokenString string) (any, error)
}

// Compile-time interface check.
var _ auth.Auth[int] = (*Auth[int])(nil)

// New returns a JWT Auth with the given extractor and validator.
func New[Req any](
	extract func(ctx context.Context, req Req) (string, error),
	validate func(ctx context.Context, tokenString string) (any, error),
) *Auth[Req] {
	return &Auth[Req]{
		ExtractToken:  extract,
		ValidateToken: validate,
	}
}

// Authenticate extracts the JWT token, validates it, and embeds the claims
// into the context.
func (a *Auth[Req]) Authenticate(ctx context.Context, req Req) (context.Context, error) {
	tokenString, err := a.ExtractToken(ctx, req)
	if err != nil {
		return ctx, fmt.Errorf("extract token: %w", err)
	}

	claims, err := a.ValidateToken(ctx, tokenString)
	if err != nil {
		return ctx, fmt.Errorf("validate token: %w", err)
	}

	return auth.SetAuth(ctx, claims), nil
}

// ErrNoToken is returned when no JWT token is found in the request.
var ErrNoToken = errors.New("no token found")
