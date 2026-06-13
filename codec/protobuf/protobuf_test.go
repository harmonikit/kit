package protobuf_test

import (
	"bytes"
	"context"
	"testing"

	protobufcodec "github.com/harmonikit/kit/codec/protobuf"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func TestCodec_Decode(t *testing.T) {
	// Use value types (not pointers) — Codec takes &req internally.
	codec := protobufcodec.NewCodec[wrapperspb.StringValue, wrapperspb.Int32Value]()

	data, _ := proto.Marshal(wrapperspb.String("hello"))
	input := bytes.NewReader(data)
	req, err := codec.Decode(context.Background(), input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if req.Value != "hello" {
		t.Fatalf("got %q, want %q", req.Value, "hello")
	}
}

func TestCodec_Encode(t *testing.T) {
	codec := protobufcodec.NewCodec[wrapperspb.StringValue, wrapperspb.Int32Value]()

	var buf bytes.Buffer
	resp := wrapperspb.Int32Value{Value: 42}
	err := codec.Encode(context.Background(), &buf, resp)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.Len() == 0 {
		t.Fatal("encoded output should not be empty")
	}
}

func TestCodec_RoundTrip(t *testing.T) {
	codec := protobufcodec.NewCodec[wrapperspb.StringValue, wrapperspb.StringValue]()

	orig := wrapperspb.StringValue{Value: "round-trip test"}
	var buf bytes.Buffer
	if err := codec.Encode(context.Background(), &buf, orig); err != nil {
		t.Fatalf("encode error: %v", err)
	}

	decoded, err := codec.Decode(context.Background(), &buf)
	if err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if decoded.Value != "round-trip test" {
		t.Fatalf("got %q, want %q", decoded.Value, "round-trip test")
	}
}

func TestCodec_Decode_Empty(t *testing.T) {
	codec := protobufcodec.NewCodec[wrapperspb.StringValue, wrapperspb.Int32Value]()

	input := bytes.NewReader([]byte{})
	req, err := codec.Decode(context.Background(), input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if req.Value != "" {
		t.Fatalf("got %q, want empty string", req.Value)
	}
}

func TestCodec_Decode_Invalid(t *testing.T) {
	codec := protobufcodec.NewCodec[wrapperspb.StringValue, wrapperspb.Int32Value]()

	input := bytes.NewReader([]byte{0xFF, 0xFE, 0xFD})
	_, err := codec.Decode(context.Background(), input)
	if err == nil {
		t.Fatal("expected decode error for invalid protobuf")
	}
}

func TestCodec_NonProtoMessage(t *testing.T) {
	type notProto struct{ Name string }
	codec := protobufcodec.NewCodec[notProto, notProto]()

	var buf bytes.Buffer
	err := codec.Encode(context.Background(), &buf, notProto{Name: "test"})
	if err == nil {
		t.Fatal("expected error for non-proto-message type")
	}
}
