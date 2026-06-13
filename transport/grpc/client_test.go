package grpc_test

import (
	"context"
	"errors"
	"testing"

	grpctransport "github.com/harmonikit/kit/transport/grpc"
)

func TestClient_EncodeError(t *testing.T) {
	wantErr := errors.New("encode failure")
	enc := func(ctx context.Context, req string) (any, error) { return nil, wantErr }
	dec := func(ctx context.Context, grpcResp any) (string, error) { return grpcResp.(string), nil }

	client := grpctransport.NewClient[string, string](nil, "/Svc/Method", enc, dec)
	ep := client.Endpoint()

	_, err := ep(context.Background(), "test")
	if err == nil {
		t.Fatal("expected error from encode failure")
	}
}
