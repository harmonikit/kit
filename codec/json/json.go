// Package json provides a JSON codec implementing transport.Codec.
//
// Example:
//
//	codec := json.NewCodec[MyReq, MyResp]()
//	req, err := codec.Decode(ctx, r)
//	err = codec.Encode(ctx, w, resp)
package json

import (
	"context"
	"encoding/json"
	"io"
)

// Codec implements transport.Codec using encoding/json.
type Codec[Req, Resp any] struct{}

// NewCodec returns a new JSON Codec.
func NewCodec[Req, Resp any]() *Codec[Req, Resp] {
	return &Codec[Req, Resp]{}
}

// Decode reads JSON from r and unmarshals it into a request.
func (c *Codec[Req, Resp]) Decode(_ context.Context, r io.Reader) (Req, error) {
	var req Req
	if err := json.NewDecoder(r).Decode(&req); err != nil {
		return req, err
	}
	return req, nil
}

// Encode marshals a response to JSON and writes it to w.
func (c *Codec[Req, Resp]) Encode(_ context.Context, w io.Writer, resp Resp) error {
	return json.NewEncoder(w).Encode(resp)
}
