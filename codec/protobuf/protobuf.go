// Package protobuf provides a Protocol Buffers codec implementing
// transport.Codec.
//
// Example:
//
//	codec := protobuf.NewCodec[MyReq, MyResp]()
//	req, err := codec.Decode(ctx, r)
//	err = codec.Encode(ctx, w, resp)
package protobuf

import (
	"context"
	"fmt"
	"io"

	"google.golang.org/protobuf/proto"
)

// Codec implements transport.Codec using Protocol Buffers encoding.
type Codec[Req, Resp any] struct{}

// NewCodec returns a new protobuf Codec.
// Req and Resp must implement proto.Message.
func NewCodec[Req, Resp any]() *Codec[Req, Resp] {
	return &Codec[Req, Resp]{}
}

// Decode reads bytes from r and unmarshals them using proto.Unmarshal.
func (c *Codec[Req, Resp]) Decode(ctx context.Context, r io.Reader) (Req, error) {
	var req Req
	data, err := io.ReadAll(r)
	if err != nil {
		return req, fmt.Errorf("read protobuf: %w", err)
	}

	msg, ok := any(&req).(proto.Message)
	if !ok {
		return req, fmt.Errorf("request type %T does not implement proto.Message", req)
	}

	if err := proto.Unmarshal(data, msg); err != nil {
		return req, fmt.Errorf("unmarshal protobuf: %w", err)
	}
	return req, nil
}

// Encode marshals a response using proto.Marshal and writes it to w.
func (c *Codec[Req, Resp]) Encode(ctx context.Context, w io.Writer, resp Resp) error {
	msg, ok := any(&resp).(proto.Message)
	if !ok {
		return fmt.Errorf("response type %T does not implement proto.Message", resp)
	}

	data, err := proto.Marshal(msg)
	if err != nil {
		return fmt.Errorf("marshal protobuf: %w", err)
	}

	_, err = w.Write(data)
	return err
}
