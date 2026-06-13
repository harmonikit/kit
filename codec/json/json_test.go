package json_test

import (
	"bytes"
	"context"
	"strings"
	"testing"

	jsoncodec "github.com/harmonikit/kit/codec/json"
)

type testReq struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type testResp struct {
	ID     int    `json:"id"`
	Status string `json:"status"`
}

func TestCodec_Decode(t *testing.T) {
	codec := jsoncodec.NewCodec[testReq, testResp]()

	input := strings.NewReader(`{"name":"alice","age":30}`)
	req, err := codec.Decode(context.Background(), input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if req.Name != "alice" {
		t.Fatalf("got name %q, want %q", req.Name, "alice")
	}
	if req.Age != 30 {
		t.Fatalf("got age %d, want %d", req.Age, 30)
	}
}

func TestCodec_Encode(t *testing.T) {
	codec := jsoncodec.NewCodec[testReq, testResp]()

	var buf bytes.Buffer
	resp := testResp{ID: 42, Status: "ok"}

	err := codec.Encode(context.Background(), &buf, resp)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := `{"id":42,"status":"ok"}` + "\n"
	if buf.String() != want {
		t.Fatalf("got %q, want %q", buf.String(), want)
	}
}

func TestCodec_RoundTrip(t *testing.T) {
	codec := jsoncodec.NewCodec[testReq, testResp]()

	// Encode a response.
	var buf bytes.Buffer
	resp := testResp{ID: 99, Status: "created"}
	if err := codec.Encode(context.Background(), &buf, resp); err != nil {
		t.Fatalf("encode error: %v", err)
	}

	// Decode it back.
	codec2 := jsoncodec.NewCodec[testResp, testResp]()
	decoded, decErr := codec2.Decode(context.Background(), &buf)
	if decErr != nil {
		t.Fatalf("decode error: %v", decErr)
	}
	if decoded.ID != 99 || decoded.Status != "created" {
		t.Fatalf("got %+v, want {99 created}", decoded)
	}
}

func TestCodec_Decode_InvalidJSON(t *testing.T) {
	codec := jsoncodec.NewCodec[testReq, testResp]()

	input := strings.NewReader(`{invalid}`)
	_, err := codec.Decode(context.Background(), input)
	if err == nil {
		t.Fatal("expected decode error for invalid JSON")
	}
}

func TestCodec_Decode_Empty(t *testing.T) {
	codec := jsoncodec.NewCodec[testReq, testResp]()

	input := strings.NewReader(`{}`)
	req, err := codec.Decode(context.Background(), input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if req.Name != "" || req.Age != 0 {
		t.Fatalf("got %+v, want zero values", req)
	}
}

func TestCodec_PrimitiveTypes(t *testing.T) {
	codec := jsoncodec.NewCodec[int, string]()

	// Decode int.
	input := strings.NewReader("42")
	req, err := codec.Decode(context.Background(), input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if req != 42 {
		t.Fatalf("got %d, want 42", req)
	}

	// Encode string.
	var buf bytes.Buffer
	err = codec.Encode(context.Background(), &buf, "hello")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.String() != `"hello"`+"\n" {
		t.Fatalf("got %q, want %q", buf.String(), `"hello"`+"\n")
	}
}

func TestCodec_Encode_NilInterface(t *testing.T) {
	// Encode a nil pointer — should produce "null".
	type wrapper struct{ Ptr *int }
	codec := jsoncodec.NewCodec[testReq, wrapper]()
	var buf bytes.Buffer
	err := codec.Encode(context.Background(), &buf, wrapper{Ptr: nil})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "null") {
		t.Fatalf("got %q, want null for nil pointer", buf.String())
	}
}
