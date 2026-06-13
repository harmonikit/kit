// Package amqp provides AMQP 0.9.1 transport bindings for harmoni endpoints.
//
// This is a stub — a production version depends on github.com/rabbitmq/amqp091-go.
package amqp

import (
	"context"

	"github.com/harmonikit/harmoni/endpoint"
)

// Server wraps a harmoni endpoint as an AMQP message handler.
type Server[Req, Resp any] struct {
	endpoint endpoint.Endpoint[Req, Resp]
	dec      DecodeRequestFunc[Req]
	enc      EncodeResponseFunc[Resp]
	queue    string
}

// DecodeRequestFunc decodes an AMQP message body into a domain request.
type DecodeRequestFunc[Req any] func(ctx context.Context, body []byte) (Req, error)

// EncodeResponseFunc encodes a domain response into an AMQP message body.
type EncodeResponseFunc[Resp any] func(ctx context.Context, resp Resp) ([]byte, error)

// NewServer creates an AMQP message handler.
func NewServer[Req, Resp any](
	ep endpoint.Endpoint[Req, Resp],
	queue string,
	dec DecodeRequestFunc[Req],
	enc EncodeResponseFunc[Resp],
) *Server[Req, Resp] {
	return &Server[Req, Resp]{
		endpoint: ep,
		dec:      dec,
		enc:      enc,
		queue:    queue,
	}
}

// HandleMsg processes an AMQP message body.
func (s *Server[Req, Resp]) HandleMsg(ctx context.Context, body []byte) ([]byte, error) {
	req, err := s.dec(ctx, body)
	if err != nil {
		return nil, err
	}

	resp, err := s.endpoint(ctx, req)
	if err != nil {
		return nil, err
	}

	return s.enc(ctx, resp)
}

// Queue returns the queue name.
func (s *Server[Req, Resp]) Queue() string { return s.queue }

// Client wraps AMQP publishing as a typed harmoni endpoint.
type Client[Req, Resp any] struct {
	enc EncodeRequestFunc[Req]
	dec DecodeResponseFunc[Resp]
}

// EncodeRequestFunc encodes a domain request into an AMQP message body.
type EncodeRequestFunc[Req any] func(ctx context.Context, req Req) ([]byte, error)

// DecodeResponseFunc decodes an AMQP message body into a domain response.
type DecodeResponseFunc[Resp any] func(ctx context.Context, body []byte) (Resp, error)

// NewClient creates an AMQP client.
func NewClient[Req, Resp any](
	enc EncodeRequestFunc[Req],
	dec DecodeResponseFunc[Resp],
) *Client[Req, Resp] {
	return &Client[Req, Resp]{enc: enc, dec: dec}
}

// Endpoint returns a harmoni endpoint that publishes via AMQP.
func (c *Client[Req, Resp]) Endpoint() endpoint.Endpoint[Req, Resp] {
	return func(ctx context.Context, req Req) (Resp, error) {
		data, err := c.enc(ctx, req)
		if err != nil {
			var zero Resp
			return zero, err
		}

		resp, err := c.dec(ctx, data)
		if err != nil {
			var zero Resp
			return zero, err
		}
		return resp, nil
	}
}
