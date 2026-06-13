// Package nats provides NATS transport bindings for harmoni endpoints.
//
// It wraps endpoints as NATS message handlers and provides client-side
// endpoint wrappers for NATS request-reply.
//
// This is a stub — a production version depends on github.com/nats-io/nats.go.
package nats

import (
	"context"

	"github.com/harmonikit/harmoni/endpoint"
)

// Server wraps a harmoni endpoint as a NATS message handler.
type Server[Req, Resp any] struct {
	endpoint endpoint.Endpoint[Req, Resp]
	dec      DecodeRequestFunc[Req]
	enc      EncodeResponseFunc[Resp]
	subject  string
}

// DecodeRequestFunc decodes a NATS message payload into a domain request.
type DecodeRequestFunc[Req any] func(ctx context.Context, data []byte) (Req, error)

// EncodeResponseFunc encodes a domain response into a NATS message payload.
type EncodeResponseFunc[Resp any] func(ctx context.Context, resp Resp) ([]byte, error)

// NewServer creates a NATS handler from an endpoint.
func NewServer[Req, Resp any](
	ep endpoint.Endpoint[Req, Resp],
	subject string,
	dec DecodeRequestFunc[Req],
	enc EncodeResponseFunc[Resp],
) *Server[Req, Resp] {
	return &Server[Req, Resp]{
		endpoint: ep,
		dec:      dec,
		enc:      enc,
		subject:  subject,
	}
}

// HandleMsg processes a NATS message.
func (s *Server[Req, Resp]) HandleMsg(ctx context.Context, data []byte) ([]byte, error) {
	req, err := s.dec(ctx, data)
	if err != nil {
		return nil, err
	}

	resp, err := s.endpoint(ctx, req)
	if err != nil {
		return nil, err
	}

	return s.enc(ctx, resp)
}

// Subject returns the NATS subject this server handles.
func (s *Server[Req, Resp]) Subject() string { return s.subject }

// Client wraps NATS request-reply as a typed harmoni endpoint.
type Client[Req, Resp any] struct {
	enc EncodeRequestFunc[Req]
	dec DecodeResponseFunc[Resp]
}

// EncodeRequestFunc encodes a domain request into a NATS message payload.
type EncodeRequestFunc[Req any] func(ctx context.Context, req Req) ([]byte, error)

// DecodeResponseFunc decodes a NATS message payload into a domain response.
type DecodeResponseFunc[Resp any] func(ctx context.Context, data []byte) (Resp, error)

// NewClient creates a NATS client.
func NewClient[Req, Resp any](
	enc EncodeRequestFunc[Req],
	dec DecodeResponseFunc[Resp],
) *Client[Req, Resp] {
	return &Client[Req, Resp]{enc: enc, dec: dec}
}

// Endpoint returns a harmoni endpoint that publishes via NATS.
func (c *Client[Req, Resp]) Endpoint() endpoint.Endpoint[Req, Resp] {
	return func(ctx context.Context, req Req) (Resp, error) {
		data, err := c.enc(ctx, req)
		if err != nil {
			var zero Resp
			return zero, err
		}

		resp, err := c.dec(ctx, data) // In production: nc.Request(subject, data, timeout)
		if err != nil {
			var zero Resp
			return zero, err
		}
		return resp, nil
	}
}
