package json_test

import (
	"bytes"
	"context"
	"testing"

	jsoncodec "github.com/harmonikit/kit/codec/json"
)

func BenchmarkCodec_Decode(b *testing.B) {
	codec := jsoncodec.NewCodec[testReq, testResp]()
	data := []byte(`{"name":"alice","age":30}`)

	for range b.N {
		_, _ = codec.Decode(context.Background(), bytes.NewReader(data))
	}
}

func BenchmarkCodec_Encode(b *testing.B) {
	codec := jsoncodec.NewCodec[testReq, testResp]()
	resp := testResp{ID: 42, Status: "ok"}

	for range b.N {
		var buf bytes.Buffer
		_ = codec.Encode(context.Background(), &buf, resp)
	}
}
