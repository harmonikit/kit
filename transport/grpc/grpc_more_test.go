package grpc_test

import (
	"context"
	"testing"

	harmoniendpoint "github.com/harmonikit/harmoni/endpoint"
	grpctransport "github.com/harmonikit/kit/transport/grpc"
	"google.golang.org/grpc"
)

func TestServer_RegisterService(t *testing.T) {
	ep := harmoniendpoint.Endpoint[string, string](func(ctx context.Context, req string) (string, error) {
		return "ok", nil
	})
	dec := func(ctx context.Context, grpcReq any) (string, error) { return grpcReq.(string), nil }
	enc := func(ctx context.Context, resp string) (any, error) { return resp, nil }

	server := grpctransport.NewServer(ep, dec, enc)
	gs := grpc.NewServer()

	// RegisterService with a simple registration function.
	server.RegisterService(gs, func(s grpc.ServiceRegistrar, h grpc.UnaryHandler) {
		// In production this would be the generated RegisterXxxServer.
		_ = s
		_ = h
	})
}

func TestServer_NewServer_WithOption(t *testing.T) {
	ep := harmoniendpoint.Endpoint[string, string](func(ctx context.Context, req string) (string, error) {
		return req, nil
	})
	dec := func(ctx context.Context, grpcReq any) (string, error) { return grpcReq.(string), nil }
	enc := func(ctx context.Context, resp string) (any, error) { return resp, nil }

	server := grpctransport.NewServer(ep, dec, enc)
	if server == nil {
		t.Fatal("expected non-nil server")
	}
}


func TestServer_WithOption(t *testing.T) {
	ep := harmoniendpoint.Endpoint[string, string](func(ctx context.Context, req string) (string, error) {
		return req, nil
	})
	dec := func(ctx context.Context, grpcReq any) (string, error) { return grpcReq.(string), nil }
	enc := func(ctx context.Context, resp string) (any, error) { return resp, nil }

	// NewServer with no options should not panic.
	server := grpctransport.NewServer(ep, dec, enc)
	_ = server.Handle()
}
