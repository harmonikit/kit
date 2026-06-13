// Package grpc provides gRPC transport bindings for harmoni endpoints.
//
// It wraps endpoints as gRPC unary handlers and provides client-side
// endpoint wrappers for calling remote gRPC services.
package grpc

import (
	"context"
	"fmt"

	"github.com/harmonikit/harmoni/endpoint"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Server wraps a harmoni endpoint as a gRPC service handler.
type Server[Req, Resp any] struct {
	endpoint endpoint.Endpoint[Req, Resp]
	dec      DecodeRequestFunc[Req]
	enc      EncodeResponseFunc[Resp]
}

// DecodeRequestFunc decodes a gRPC request into a domain request.
type DecodeRequestFunc[Req any] func(ctx context.Context, grpcReq any) (Req, error)

// EncodeResponseFunc encodes a domain response into a gRPC response.
type EncodeResponseFunc[Resp any] func(ctx context.Context, resp Resp) (any, error)

// ServerOption configures a Server.
type ServerOption[Req, Resp any] func(*Server[Req, Resp])

// NewServer creates a gRPC handler from an endpoint and codec functions.
func NewServer[Req, Resp any](
	ep endpoint.Endpoint[Req, Resp],
	dec DecodeRequestFunc[Req],
	enc EncodeResponseFunc[Resp],
	opts ...ServerOption[Req, Resp],
) *Server[Req, Resp] {
	s := &Server[Req, Resp]{
		endpoint: ep,
		dec:      dec,
		enc:      enc,
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

// Handle returns a gRPC unary handler function.
func (s *Server[Req, Resp]) Handle() grpc.UnaryHandler {
	return func(ctx context.Context, grpcReq any) (any, error) {
		req, err := s.dec(ctx, grpcReq)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("decode request: %v", err))
		}

		resp, err := s.endpoint(ctx, req)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "endpoint error: %v", err)
		}

		grpcResp, err := s.enc(ctx, resp)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "encode response: %v", err)
		}
		return grpcResp, nil
	}
}

// RegisterService registers the handler on a gRPC server using the provided
// registration function. The caller must provide the generated RegisterXxxServer
// function from the protobuf service definition.
func (s *Server[Req, Resp]) RegisterService(grpcServer *grpc.Server, register func(grpc.ServiceRegistrar, grpc.UnaryHandler)) {
	register(grpcServer, s.Handle())
}
