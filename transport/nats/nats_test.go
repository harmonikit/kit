package nats_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/harmonikit/harmoni/endpoint"
	nats "github.com/harmonikit/kit/transport/nats"
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

	server := nats.NewServer(ep, "test.subject", dec, enc)
	if server.Subject() != "test.subject" {
		t.Fatalf("got %q, want test.subject", server.Subject())
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

func TestServer_DecodeError(t *testing.T) {
	ep := endpoint.Endpoint[int, int](func(ctx context.Context, req int) (int, error) {
		return req, nil
	})

	dec := func(ctx context.Context, data []byte) (int, error) {
		return 0, errors.New("bad data")
	}

	enc := func(ctx context.Context, resp int) ([]byte, error) {
		return nil, nil
	}

	server := nats.NewServer(ep, "test", dec, enc)
	_, err := server.HandleMsg(context.Background(), []byte("bad"))
	if err == nil {
		t.Fatal("expected decode error")
	}
}
