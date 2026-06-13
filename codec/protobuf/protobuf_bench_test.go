package protobuf_test

import (
	"bytes"
	"context"
	"testing"

	protobufcodec "github.com/harmonikit/kit/codec/protobuf"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func BenchmarkCodec_Decode(b *testing.B) {
	codec := protobufcodec.NewCodec[wrapperspb.StringValue, wrapperspb.Int32Value]()
	data, _ := proto.Marshal(wrapperspb.String("benchmark"))

	for range b.N {
		_, _ = codec.Decode(context.Background(), bytes.NewReader(data))
	}
}

func BenchmarkCodec_Encode(b *testing.B) {
	codec := protobufcodec.NewCodec[wrapperspb.StringValue, wrapperspb.Int32Value]()
	resp := wrapperspb.Int32Value{Value: 42}

	for range b.N {
		var buf bytes.Buffer
		_ = codec.Encode(context.Background(), &buf, resp)
	}
}
