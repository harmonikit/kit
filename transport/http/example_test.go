package http_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"

	harmoniendpoint "github.com/harmonikit/harmoni/endpoint"
	httptransport "github.com/harmonikit/kit/transport/http"
)

func ExampleServer() {
	// Define a typed endpoint.
	add := harmoniendpoint.Endpoint[int, int](func(ctx context.Context, req int) (int, error) {
		return req + 1, nil
	})

	// Decode: read JSON from request body.
	decode := func(ctx context.Context, r *http.Request) (int, error) {
		body, _ := io.ReadAll(r.Body)
		var v int
		json.Unmarshal(body, &v)
		return v, nil
	}

	// Encode: write JSON to response.
	encode := func(ctx context.Context, w http.ResponseWriter, resp int) error {
		return json.NewEncoder(w).Encode(resp)
	}

	server := httptransport.NewServer(add, decode, encode)

	// Test it with httptest.
	req := httptest.NewRequest("POST", "/", strings.NewReader("41"))
	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)

	var result int
	json.NewDecoder(rec.Body).Decode(&result)
	fmt.Println(result)
	// Output: 42
}

func ExampleServer_customErrorEncoder() {
	ep := harmoniendpoint.Endpoint[int, int](func(ctx context.Context, req int) (int, error) {
		return 0, fmt.Errorf("not found: %d", req)
	})

	decode := func(ctx context.Context, r *http.Request) (int, error) { return 1, nil }
	encode := func(ctx context.Context, w http.ResponseWriter, resp int) error { return nil }

	customErr := func(ctx context.Context, w http.ResponseWriter, err error) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, `{"error":%q}`, err.Error())
	}

	server := httptransport.NewServer(ep, decode, encode,
		httptransport.WithErrorEncoder[int, int](customErr))

	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)

	fmt.Printf("status=%d body=%s", rec.Code, rec.Body.String())
	// Output: status=404 body={"error":"not found: 1"}
}
