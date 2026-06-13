// Package fiber provides a Fiber transport binding for harmoni endpoints.
//
// It wraps a typed endpoint as a Fiber handler, giving access to Fiber's
// performance, middleware ecosystem, and ergonomics while keeping the
// harmoni type-safe endpoint abstraction.
//
// Example:
//
//	app := fiber.New()
//	server := fiber.NewServer(myEndpoint, decodeRequest, encodeResponse)
//	app.Post("/users", server.Handle)
//	app.Listen(":8080")
package fiber

import (
	"context"

	harmoniendpoint "github.com/harmonikit/harmoni/endpoint"
	"github.com/gofiber/fiber/v3"
)

// DecodeRequestFunc decodes a Fiber request into a domain request.
type DecodeRequestFunc[Req any] func(ctx context.Context, c fiber.Ctx) (Req, error)

// EncodeResponseFunc encodes a domain response into a Fiber response.
type EncodeResponseFunc[Resp any] func(ctx context.Context, c fiber.Ctx, resp Resp) error

// EncodeErrorFunc encodes an error into a Fiber response.
type EncodeErrorFunc func(ctx context.Context, c fiber.Ctx, err error) error

// Server wraps a harmoni endpoint as a Fiber handler.
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

// defaultErrorEncoder returns a 500 status with the error message.
func defaultErrorEncoder(ctx context.Context, c fiber.Ctx, err error) error {
	return c.Status(500).SendString(err.Error())
}

// NewServer creates a Fiber handler from an endpoint, request decoder, and
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

// Handle returns a Fiber handler. Register it with app.Get, app.Post, etc.
func (s *Server[Req, Resp]) Handle(c fiber.Ctx) error {
	ctx := c.Context()

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
