package http_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	httptransport "github.com/harmonikit/kit/transport/http"
)

func TestClient_Endpoint(t *testing.T) {
	// Start a test server that echoes back the request body as the response.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var v int
		json.Unmarshal(body, &v)
		json.NewEncoder(w).Encode(v * 2)
	}))
	defer ts.Close()

	u, err := url.Parse(ts.URL)
	if err != nil {
		t.Fatalf("parse url: %v", err)
	}

	enc := func(ctx context.Context, req int) (io.Reader, error) {
		data, _ := json.Marshal(req)
		return bytes.NewReader(data), nil
	}

	dec := func(ctx context.Context, r *http.Response) (int, error) {
		var v int
		json.NewDecoder(r.Body).Decode(&v)
		return v, nil
	}

	client := httptransport.NewClient("POST", u, enc, dec)
	ep := client.Endpoint()

	resp, err := ep(context.Background(), 21)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp != 42 {
		t.Fatalf("got %d, want 42", resp)
	}
}

func TestClient_ErrorStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "not found", http.StatusNotFound)
	}))
	defer ts.Close()

	u, _ := url.Parse(ts.URL)

	enc := func(ctx context.Context, req string) (io.Reader, error) {
		return strings.NewReader(req), nil
	}

	dec := func(ctx context.Context, r *http.Response) (string, error) {
		body, _ := io.ReadAll(r.Body)
		return string(body), nil
	}

	client := httptransport.NewClient("GET", u, enc, dec)
	ep := client.Endpoint()

	_, err := ep(context.Background(), "hello")
	if err == nil {
		t.Fatal("expected error for non-2xx status")
	}
}
