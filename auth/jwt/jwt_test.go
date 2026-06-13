package jwt_test

import (
	"context"
	"errors"
	"testing"

	"github.com/harmonikit/harmoni/auth"
	jwtauth "github.com/harmonikit/kit/auth/jwt"
)

func TestAuth_Authenticate(t *testing.T) {
	extract := func(ctx context.Context, req int) (string, error) {
		return "token-123", nil
	}
	validate := func(ctx context.Context, tokenString string) (any, error) {
		if tokenString != "token-123" {
			return nil, errors.New("invalid token")
		}
		return map[string]any{"sub": "user-42"}, nil
	}

	ja := jwtauth.New[int](extract, validate)
	ctx, err := ja.Authenticate(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	claims, ok := auth.GetAuth(ctx)
	if !ok {
		t.Fatal("expected claims in context")
	}
	m := claims.(map[string]any)
	if m["sub"] != "user-42" {
		t.Fatalf("got %v, want user-42", m["sub"])
	}
}

func TestAuth_ExtractError(t *testing.T) {
	wantErr := errors.New("no token")
	extract := func(ctx context.Context, req int) (string, error) {
		return "", wantErr
	}
	validate := func(ctx context.Context, tokenString string) (any, error) {
		return nil, nil
	}

	ja := jwtauth.New[int](extract, validate)
	_, err := ja.Authenticate(context.Background(), 1)
	if !errors.Is(err, wantErr) {
		t.Fatalf("got %v, want %v", err, wantErr)
	}
}

func TestAuth_ValidateError(t *testing.T) {
	wantErr := errors.New("token expired")
	extract := func(ctx context.Context, req int) (string, error) {
		return "bad-token", nil
	}
	validate := func(ctx context.Context, tokenString string) (any, error) {
		return nil, wantErr
	}

	ja := jwtauth.New[int](extract, validate)
	_, err := ja.Authenticate(context.Background(), 1)
	if !errors.Is(err, wantErr) {
		t.Fatalf("got %v, want %v", err, wantErr)
	}
}

func TestAuth_Interface(t *testing.T) {
	var _ auth.Auth[int] = jwtauth.New[int](nil, nil)
}
