// Package msgpack provides a MessagePack codec implementing transport.Codec.
//
// This is a stub — a production version depends on a MessagePack library
// like github.com/vmihailenco/msgpack/v5 or github.com/ugorji/go/codec.
package msgpack

import (
	"context"
	"encoding/json"
	"io"
)

// Codec implements transport.Codec using JSON as a stand-in for MessagePack.
// In production, replace encoding/json with a MessagePack encoder.
type Codec[Req, Resp any] struct{}

// NewCodec returns a new Codec.
func NewCodec[Req, Resp any]() *Codec[Req, Resp] {
	return &Codec[Req, Resp]{}
}

// Decode reads JSON from r. Replace with msgpack in production.
func (c *Codec[Req, Resp]) Decode(_ context.Context, r io.Reader) (Req, error) {
	var req Req
	if err := json.NewDecoder(r).Decode(&req); err != nil {
		return req, err
	}
	return req, nil
}

// Encode writes JSON to w. Replace with msgpack in production.
func (c *Codec[Req, Resp]) Encode(_ context.Context, w io.Writer, resp Resp) error {
	return json.NewEncoder(w).Encode(resp)
}
