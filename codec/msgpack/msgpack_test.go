package msgpack_test

import (
	"bytes"
	"context"
	"strings"
	"testing"

	msgpack "github.com/harmonikit/kit/codec/msgpack"
)

func TestCodec_Decode(t *testing.T) {
	codec := msgpack.NewCodec[int, string]()

	input := strings.NewReader("42")
	req, err := codec.Decode(context.Background(), input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if req != 42 {
		t.Fatalf("got %d, want 42", req)
	}
}

func TestCodec_Encode(t *testing.T) {
	codec := msgpack.NewCodec[int, string]()

	var buf bytes.Buffer
	err := codec.Encode(context.Background(), &buf, "hello")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.String() != `"hello"`+"\n" {
		t.Fatalf("got %q, want %q", buf.String(), `"hello"`+"\n")
	}
}

func TestCodec_RoundTrip(t *testing.T) {
	codec := msgpack.NewCodec[int, int]()

	var buf bytes.Buffer
	if err := codec.Encode(context.Background(), &buf, 42); err != nil {
		t.Fatalf("encode error: %v", err)
	}

	val, err := codec.Decode(context.Background(), &buf)
	if err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if val != 42 {
		t.Fatalf("got %d, want 42", val)
	}
}

func TestCodec_Decode_Invalid(t *testing.T) {
	codec := msgpack.NewCodec[int, int]()
	input := strings.NewReader("not-json")
	_, err := codec.Decode(context.Background(), input)
	if err == nil {
		t.Fatal("expected decode error for invalid data")
	}
}
