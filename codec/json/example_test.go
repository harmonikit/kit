package json_test

import (
	"bytes"
	"context"
	"fmt"

	jsoncodec "github.com/harmonikit/kit/codec/json"
)

func ExampleCodec() {
	codec := jsoncodec.NewCodec[int, string]()

	// Decode a request.
	req, _ := codec.Decode(context.Background(), bytes.NewReader([]byte("42")))
	fmt.Println(req)

	// Encode a response.
	var buf bytes.Buffer
	codec.Encode(context.Background(), &buf, "hello")
	fmt.Print(buf.String())
	// Output:
	// 42
	// "hello"
}
