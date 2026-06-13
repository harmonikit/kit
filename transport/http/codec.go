package http

import (
	"context"
	"net/http"
)

// DecodeRequestFunc decodes an HTTP request into a domain request.
type DecodeRequestFunc[Req any] func(ctx context.Context, r *http.Request) (Req, error)

// EncodeResponseFunc encodes a domain response into an HTTP response.
type EncodeResponseFunc[Resp any] func(ctx context.Context, w http.ResponseWriter, resp Resp) error

// EncodeErrorFunc encodes an error into an HTTP response.
type EncodeErrorFunc func(ctx context.Context, w http.ResponseWriter, err error)
