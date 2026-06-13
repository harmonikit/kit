// Package legacy provides an adapter from the legacy go-kit/kit endpoint.Endpoint
// (the interface{} version) to the modern harmoni endpoint.Endpoint[any, any].
//
// This enables incremental migration from go-kit to harmonikit.
//
// Example:
//
//	import gokitendpoint "github.com/go-kit/kit/endpoint"
//	import "github.com/harmonikit/kit/legacy"
//
//	oldEP := gokitendpoint.Endpoint(func(ctx context.Context, request interface{}) (interface{}, error) {
//	    return request, nil
//	})
//	newEP := legacy.Adapt(oldEP)
//	// newEP is now a harmoni endpoint.Endpoint[any, any]
package legacy

import (
	"context"
	"fmt"

	"github.com/harmonikit/harmoni/endpoint"
)

// Endpoint is the legacy go-kit endpoint type (interface{} based).
type Endpoint func(ctx context.Context, request any) (any, error)

// Adapt wraps a legacy go-kit endpoint as a modern harmoni endpoint.
// Both request and response are any.
func Adapt(legacy Endpoint) endpoint.Endpoint[any, any] {
	return func(ctx context.Context, req any) (any, error) {
		return legacy(ctx, req)
	}
}

// AdaptTyped wraps a legacy go-kit endpoint as a modern typed harmoni endpoint.
// The caller asserts the types at the boundary; if the assertion fails at
// runtime, the error is returned.
func AdaptTyped[Req, Resp any](legacy Endpoint) endpoint.Endpoint[Req, Resp] {
	return func(ctx context.Context, req Req) (Resp, error) {
		result, err := legacy(ctx, req)
		if err != nil {
			var zero Resp
			return zero, err
		}
		resp, ok := result.(Resp)
		if !ok {
			var zero Resp
			return zero, fmt.Errorf("legacy: unexpected response type: got %T, want %T", result, zero)
		}
		return resp, nil
	}
}
