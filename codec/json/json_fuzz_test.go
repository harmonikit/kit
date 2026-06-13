package json_test

import (
	"bytes"
	"context"
	"testing"

	jsoncodec "github.com/harmonikit/kit/codec/json"
)

func FuzzCodec_Decode(f *testing.F) {
	f.Add([]byte(`{"name":"alice","age":30}`))
	f.Add([]byte(`{}`))
	f.Add([]byte(`{"name":""}`))
	f.Add([]byte(`invalid`))

	f.Fuzz(func(t *testing.T, data []byte) {
		codec := jsoncodec.NewCodec[testReq, testResp]()
		req, err := codec.Decode(context.Background(), bytes.NewReader(data))
		if err != nil {
			// Invalid JSON is expected — just make sure we don't panic.
			return
		}
		// If decode succeeded, re-encoding should also work.
		// (We can't re-encode testReq as testResp, so just verify no panic.)
		_ = req
	})
}

func FuzzCodec_Encode(f *testing.F) {
	f.Add(int64(0), "ok")
	f.Add(int64(42), "created")

	f.Fuzz(func(t *testing.T, id int64, status string) {
		codec := jsoncodec.NewCodec[testReq, testResp]()
		var buf bytes.Buffer
		resp := testResp{ID: int(id), Status: status}
		err := codec.Encode(context.Background(), &buf, resp)
		if err != nil {
			t.Fatalf("encode should never fail: %v", err)
		}
		if buf.Len() == 0 {
			t.Fatal("encoded output should not be empty")
		}
	})
}
