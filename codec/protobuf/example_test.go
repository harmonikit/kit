package protobuf_test

import (
	"bytes"
	"context"
	"fmt"

	protobufcodec "github.com/harmonikit/kit/codec/protobuf"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func ExampleCodec() {
	codec := protobufcodec.NewCodec[wrapperspb.StringValue, wrapperspb.Int32Value]()

	// Decode a request.
	data, _ := proto.Marshal(wrapperspb.String("hello"))
	req, _ := codec.Decode(context.Background(), bytes.NewReader(data))
	fmt.Println(req.Value)

	// Encode a response.
	var buf bytes.Buffer
	codec.Encode(context.Background(), &buf, wrapperspb.Int32Value{Value: 42})
	fmt.Println(buf.Len() > 0)
	// Output:
	// hello
	// true
}
