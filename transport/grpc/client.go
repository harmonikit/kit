package grpc

import (
	"context"
	"fmt"

	"github.com/harmonikit/harmoni/endpoint"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Client wraps a gRPC connection as a typed harmoni endpoint.
type Client[Req, Resp any] struct {
	conn   *grpc.ClientConn
	method string
	enc    EncodeRequestFunc[Req]
	dec    DecodeResponseFunc[Resp]
}

// EncodeRequestFunc encodes a domain request into a gRPC request message.
type EncodeRequestFunc[Req any] func(ctx context.Context, req Req) (any, error)

// DecodeResponseFunc decodes a gRPC response message into a domain response.
type DecodeResponseFunc[Resp any] func(ctx context.Context, grpcResp any) (Resp, error)

// NewClient creates a gRPC client.
func NewClient[Req, Resp any](
	conn *grpc.ClientConn,
	method string,
	enc EncodeRequestFunc[Req],
	dec DecodeResponseFunc[Resp],
) *Client[Req, Resp] {
	return &Client[Req, Resp]{
		conn:   conn,
		method: method,
		enc:    enc,
		dec:    dec,
	}
}

// Endpoint returns a harmoni endpoint that calls the gRPC method.
func (c *Client[Req, Resp]) Endpoint() endpoint.Endpoint[Req, Resp] {
	return func(ctx context.Context, req Req) (Resp, error) {
		grpcReq, err := c.enc(ctx, req)
		if err != nil {
			var zero Resp
			return zero, fmt.Errorf("encode request: %w", err)
		}

		var grpcResp any // Will be filled by the codec.
		err = c.conn.Invoke(ctx, c.method, grpcReq, &grpcResp)
		if err != nil {
			if st, ok := status.FromError(err); ok && st.Code() == codes.Unimplemented {
				var zero Resp
				return zero, fmt.Errorf("grpc unimplemented: %s", c.method)
			}
			var zero Resp
			return zero, fmt.Errorf("grpc invoke: %w", err)
		}

		resp, err := c.dec(ctx, grpcResp)
		if err != nil {
			var zero Resp
			return zero, fmt.Errorf("decode response: %w", err)
		}
		return resp, nil
	}
}
