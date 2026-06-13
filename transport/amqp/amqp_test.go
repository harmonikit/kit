package amqp_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/harmonikit/harmoni/endpoint"
	amqp "github.com/harmonikit/kit/transport/amqp"
)

func TestServer_HandleMsg(t *testing.T) {
	ep := endpoint.Endpoint[int, int](func(ctx context.Context, req int) (int, error) {
		return req * 2, nil
	})

	dec := func(ctx context.Context, data []byte) (int, error) {
		var v int
		json.Unmarshal(data, &v)
		return v, nil
	}

	enc := func(ctx context.Context, resp int) ([]byte, error) {
		return json.Marshal(resp)
	}

	server := amqp.NewServer(ep, "test.queue", dec, enc)
	if server.Queue() != "test.queue" {
		t.Fatalf("got %q, want test.queue", server.Queue())
	}

	resp, err := server.HandleMsg(context.Background(), []byte("21"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var result int
	json.Unmarshal(resp, &result)
	if result != 42 {
		t.Fatalf("got %d, want 42", result)
	}
}

func TestServer_HandleMsg_DecodeError(t *testing.T) {
	ep := endpoint.Endpoint[int, int](func(ctx context.Context, req int) (int, error) {
		return req, nil
	})

	dec := func(ctx context.Context, data []byte) (int, error) {
		return 0, errors.New("bad data")
	}

	enc := func(ctx context.Context, resp int) ([]byte, error) { return nil, nil }

	server := amqp.NewServer(ep, "test", dec, enc)
	_, err := server.HandleMsg(context.Background(), []byte("bad"))
	if err == nil {
		t.Fatal("expected decode error")
	}
}
