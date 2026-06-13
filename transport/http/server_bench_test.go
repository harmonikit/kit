package http_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	harmoniendpoint "github.com/harmonikit/harmoni/endpoint"
	httptransport "github.com/harmonikit/kit/transport/http"
)

func BenchmarkServer_ServeHTTP(b *testing.B) {
	ep := harmoniendpoint.Endpoint[int, int](func(ctx context.Context, req int) (int, error) {
		return req * 2, nil
	})

	dec := func(ctx context.Context, r *http.Request) (int, error) {
		body, _ := io.ReadAll(r.Body)
		var v int
		json.Unmarshal(body, &v)
		return v, nil
	}

	enc := func(ctx context.Context, w http.ResponseWriter, resp int) error {
		return json.NewEncoder(w).Encode(resp)
	}

	server := httptransport.NewServer(ep, dec, enc)

	for range b.N {
		body := strings.NewReader("21")
		req := httptest.NewRequest("POST", "/", body)
		rec := httptest.NewRecorder()
		server.ServeHTTP(rec, req)
	}
}
