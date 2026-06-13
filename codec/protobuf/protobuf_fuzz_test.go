package protobuf_test

import (
	"bytes"
	"context"
	"testing"

	protobufcodec "github.com/harmonikit/kit/codec/protobuf"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func FuzzCodec_Decode(f *testing.F) {
	valid, _ := proto.Marshal(wrapperspb.String("hello"))
	f.Add(valid)
	f.Add([]byte{})
	f.Add([]byte{0xFF, 0xFE})

	f.Fuzz(func(t *testing.T, data []byte) {
		codec := protobufcodec.NewCodec[wrapperspb.StringValue, wrapperspb.Int32Value]()
		req, err := codec.Decode(context.Background(), bytes.NewReader(data))
		if err != nil {
			return
		}
		_ = req
	})
}
