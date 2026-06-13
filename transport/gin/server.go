// Package gin provides a Gin transport binding for harmoni endpoints.
//
// It wraps a typed endpoint as a Gin handler, giving access to Gin's
// router, middleware ecosystem, and ergonomics while keeping the harmoni
// type-safe endpoint abstraction.
//
// Example:
//
//	r := gin.Default()
//	server := gin.NewServer(myEndpoint, decodeRequest, encodeResponse)
//	r.POST("/users", server.Handle)
//	r.Run(":8080")
package gin

import (
	"context"

	"github.com/gin-gonic/gin"
	harmoniendpoint "github.com/harmonikit/harmoni/endpoint"
)

// DecodeRequestFunc decodes a Gin request into a domain request.
type DecodeRequestFunc[Req any] func(ctx context.Context, c *gin.Context) (Req, error)

// EncodeResponseFunc encodes a domain response into a Gin response.
type EncodeResponseFunc[Resp any] func(ctx context.Context, c *gin.Context, resp Resp) error

// EncodeErrorFunc encodes an error into a Gin response.
type EncodeErrorFunc func(ctx context.Context, c *gin.Context, err error) error

// Server wraps a harmoni endpoint as a Gin handler.
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
func defaultErrorEncoder(_ context.Context, c *gin.Context, err error) error {
	c.String(500, err.Error())
	return nil
}

// NewServer creates a Gin handler from an endpoint, request decoder, and
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

// Handle returns a gin.HandlerFunc. Register it with r.GET, r.POST, etc.
func (s *Server[Req, Resp]) Handle(c *gin.Context) {
	ctx := c.Request.Context()

	req, err := s.dec(ctx, c)
	if err != nil {
		_ = s.encError(ctx, c, err)
		return
	}

	resp, err := s.endpoint(ctx, req)
	if err != nil {
		_ = s.encError(ctx, c, err)
		return
	}

	if err := s.enc(ctx, c, resp); err != nil {
		_ = s.encError(ctx, c, err)
		return
	}
}
