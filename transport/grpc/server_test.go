package grpc_test

import (
	"context"
	"errors"
	"testing"

	harmoniendpoint "github.com/harmonikit/harmoni/endpoint"
	grpctransport "github.com/harmonikit/kit/transport/grpc"
)

func TestServer_Handle_Success(t *testing.T) {
	ep := harmoniendpoint.Endpoint[string, string](func(ctx context.Context, req string) (string, error) {
		return "echo: " + req, nil
	})

	dec := func(ctx context.Context, grpcReq any) (string, error) {
		return grpcReq.(string), nil
	}

	enc := func(ctx context.Context, resp string) (any, error) {
		return resp, nil
	}

	server := grpctransport.NewServer(ep, dec, enc)
	handler := server.Handle()

	resp, err := handler(context.Background(), "hello")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp != "echo: hello" {
		t.Fatalf("got %q, want %q", resp, "echo: hello")
	}
}

func TestServer_Handle_DecodeError(t *testing.T) {
	ep := harmoniendpoint.Endpoint[string, string](func(ctx context.Context, req string) (string, error) {
		return req, nil
	})

	wantErr := errors.New("invalid request")
	dec := func(ctx context.Context, grpcReq any) (string, error) {
		return "", wantErr
	}

	enc := func(ctx context.Context, resp string) (any, error) { return resp, nil }

	server := grpctransport.NewServer(ep, dec, enc)
	handler := server.Handle()

	_, err := handler(context.Background(), struct{}{})
	if err == nil {
		t.Fatal("expected error from decode failure")
	}
}

func TestServer_Handle_EndpointError(t *testing.T) {
	wantErr := errors.New("business logic error")
	ep := harmoniendpoint.Endpoint[string, string](func(ctx context.Context, req string) (string, error) {
		return "", wantErr
	})

	dec := func(ctx context.Context, grpcReq any) (string, error) { return grpcReq.(string), nil }
	enc := func(ctx context.Context, resp string) (any, error) { return resp, nil }

	server := grpctransport.NewServer(ep, dec, enc)
	handler := server.Handle()

	_, err := handler(context.Background(), "test")
	if err == nil {
		t.Fatal("expected error from endpoint failure")
	}
}

func TestServer_Handle_EncodeError(t *testing.T) {
	ep := harmoniendpoint.Endpoint[string, string](func(ctx context.Context, req string) (string, error) {
		return "ok", nil
	})

	dec := func(ctx context.Context, grpcReq any) (string, error) { return grpcReq.(string), nil }
	wantErr := errors.New("encode failure")
	enc := func(ctx context.Context, resp string) (any, error) { return nil, wantErr }

	server := grpctransport.NewServer(ep, dec, enc)
	handler := server.Handle()

	_, err := handler(context.Background(), "test")
	if err == nil {
		t.Fatal("expected error from encode failure")
	}
}

func TestClient_Endpoint_EncodeError(t *testing.T) {
	wantErr := errors.New("encode failure")
	enc := func(ctx context.Context, req string) (any, error) { return nil, wantErr }
	dec := func(ctx context.Context, grpcResp any) (string, error) { return grpcResp.(string), nil }

	// Client with nil conn — encode fails before connection is used.
	client := grpctransport.NewClient[string, string](nil, "/Service/Method", enc, dec)
	ep := client.Endpoint()

	_, err := ep(context.Background(), "test")
	if err == nil {
		t.Fatal("expected error from encode failure")
	}
}
