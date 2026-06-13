package grpc_test

import (
	"context"
	"fmt"

	harmoniendpoint "github.com/harmonikit/harmoni/endpoint"
	grpctransport "github.com/harmonikit/kit/transport/grpc"
)

func ExampleServer() {
	ep := harmoniendpoint.Endpoint[string, string](func(ctx context.Context, req string) (string, error) {
		return "hello, " + req, nil
	})

	dec := func(ctx context.Context, grpcReq any) (string, error) {
		return grpcReq.(string), nil
	}

	enc := func(ctx context.Context, resp string) (any, error) {
		return resp, nil
	}

	server := grpctransport.NewServer(ep, dec, enc)
	handler := server.Handle()

	resp, _ := handler(context.Background(), "world")
	fmt.Println(resp)
	// Output: hello, world
}
