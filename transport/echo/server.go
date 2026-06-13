// Package echo provides an Echo transport binding for harmoni endpoints.
//
// It wraps a typed endpoint as an Echo handler, giving access to Echo's
// router, middleware ecosystem, and ergonomics while keeping the harmoni
// type-safe endpoint abstraction.
//
// Example:
//
//	e := echo.New()
//	server := echo.NewServer(myEndpoint, decodeRequest, encodeResponse)
//	e.POST("/users", server.Handle)
//	e.Start(":8080")
package echo

import (
	"context"

	"github.com/labstack/echo/v4"

	harmoniendpoint "github.com/harmonikit/harmoni/endpoint"
)

// DecodeRequestFunc decodes an Echo request into a domain request.
type DecodeRequestFunc[Req any] func(ctx context.Context, c echo.Context) (Req, error)

// EncodeResponseFunc encodes a domain response into an Echo response.
type EncodeResponseFunc[Resp any] func(ctx context.Context, c echo.Context, resp Resp) error

// EncodeErrorFunc encodes an error into an Echo response.
type EncodeErrorFunc func(ctx context.Context, c echo.Context, err error) error

// Server wraps a harmoni endpoint as an Echo handler.
type Server[Req, Resp any] struct {
	endpoint harmoniendpoint.Endpoint[Req, Resp]
	dec      DecodeRequestFunc[Req]
	enc      EncodeResponseFunc[Resp]
	encError EncodeErrorFunc
}

// ServerOption configures a Server.
type ServerOption[Req, Resp any] func(*Server[Req, Resp])

// WithErrorEncoder sets a custom error encoder.
func WithErrorEncoder[Req, Resp any](fn EncodeErrorFunc) ServerOption[Req, Resp] {
	return func(s *Server[Req, Resp]) { s.encError = fn }
}

// defaultErrorEncoder returns a 500 status with the error message via Echo.
func defaultErrorEncoder(_ context.Context, c echo.Context, err error) error {
	return c.String(500, err.Error())
}

// NewServer creates an Echo handler from an endpoint, request decoder, and
// response encoder.
func NewServer[Req, Resp any](
	ep harmoniendpoint.Endpoint[Req, Resp],
	dec DecodeRequestFunc[Req],
	enc EncodeResponseFunc[Resp],
	opts ...ServerOption[Req, Resp],
) *Server[Req, Resp] {
	s := &Server[Req, Resp]{
		endpoint: ep,
		dec:      dec,
		enc:      enc,
		encError: defaultErrorEncoder,
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

// Handle returns an Echo handler function. Register it with e.GET, e.POST, etc.
func (s *Server[Req, Resp]) Handle(c echo.Context) error {
	ctx := c.Request().Context()

	req, err := s.dec(ctx, c)
	if err != nil {
		return s.encError(ctx, c, err)
	}

	resp, err := s.endpoint(ctx, req)
	if err != nil {
		return s.encError(ctx, c, err)
	}

	return s.enc(ctx, c, resp)
}
