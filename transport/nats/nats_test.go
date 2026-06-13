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

func TestServer_EndpointError(t *testing.T) {
	wantErr := errors.New("handler error")
	ep := endpoint.Endpoint[int, int](func(ctx context.Context, req int) (int, error) {
		return 0, wantErr
	})

	dec := func(ctx context.Context, data []byte) (int, error) { return 1, nil }
	enc := func(ctx context.Context, resp int) ([]byte, error) { return nil, nil }

	server := nats.NewServer(ep, "test", dec, enc)
	_, err := server.HandleMsg(context.Background(), nil)
	if !errors.Is(err, wantErr) {
		t.Fatalf("got %v, want %v", err, wantErr)
	}
}

func TestClient(t *testing.T) {
	enc := func(ctx context.Context, req int) ([]byte, error) {
		return json.Marshal(req)
	}
	dec := func(ctx context.Context, data []byte) (int, error) {
		var v int
		json.Unmarshal(data, &v)
		return v, nil
	}

	client := nats.NewClient[int, int](enc, dec)
	ep := client.Endpoint()

	resp, err := ep(context.Background(), 42)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp != 42 {
		t.Fatalf("got %d, want 42", resp)
	}
}

func TestClient_EncodeError(t *testing.T) {
	wantErr := errors.New("encode fail")
	enc := func(ctx context.Context, req int) ([]byte, error) { return nil, wantErr }
	dec := func(ctx context.Context, data []byte) (int, error) { return 1, nil }

	client := nats.NewClient[int, int](enc, dec)
	ep := client.Endpoint()
	_, err := ep(context.Background(), 1)
	if !errors.Is(err, wantErr) {
		t.Fatalf("got %v, want %v", err, wantErr)
	}
}

func TestClient_DecodeError(t *testing.T) {
	wantErr := errors.New("decode fail")
	enc := func(ctx context.Context, req int) ([]byte, error) { return []byte("x"), nil }
	dec := func(ctx context.Context, data []byte) (int, error) { return 0, wantErr }

	client := nats.NewClient[int, int](enc, dec)
	ep := client.Endpoint()
	_, err := ep(context.Background(), 1)
	if !errors.Is(err, wantErr) {
		t.Fatalf("got %v, want %v", err, wantErr)
	}
}
