package http_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	harmoniendpoint "github.com/harmonikit/harmoni/endpoint"
	httptransport "github.com/harmonikit/kit/transport/http"
)

func TestServer_Success(t *testing.T) {
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

	body := strings.NewReader("21")
	req := httptest.NewRequest("POST", "/", body)
	rec := httptest.NewRecorder()

	server.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("got status %d, want %d", rec.Code, http.StatusOK)
	}

	var resp int
	json.NewDecoder(rec.Body).Decode(&resp)
	if resp != 42 {
		t.Fatalf("got %d, want 42", resp)
	}
}

func TestServer_DecodeError(t *testing.T) {
	ep := harmoniendpoint.Endpoint[int, int](func(ctx context.Context, req int) (int, error) {
		return req, nil
	})

	dec := func(ctx context.Context, r *http.Request) (int, error) {
		return 0, errors.New("bad request")
	}

	enc := func(ctx context.Context, w http.ResponseWriter, resp int) error {
		return json.NewEncoder(w).Encode(resp)
	}

	server := httptransport.NewServer(ep, dec, enc)

	req := httptest.NewRequest("POST", "/", strings.NewReader("not-json"))
	rec := httptest.NewRecorder()

	server.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("got status %d, want %d", rec.Code, http.StatusInternalServerError)
	}
}

func TestServer_EndpointError(t *testing.T) {
	wantErr := errors.New("business error")
	ep := harmoniendpoint.Endpoint[int, int](func(ctx context.Context, req int) (int, error) {
		return 0, wantErr
	})

	dec := func(ctx context.Context, r *http.Request) (int, error) {
		return 1, nil
	}

	enc := func(ctx context.Context, w http.ResponseWriter, resp int) error {
		return json.NewEncoder(w).Encode(resp)
	}

	server := httptransport.NewServer(ep, dec, enc)

	req := httptest.NewRequest("POST", "/", strings.NewReader("1"))
	rec := httptest.NewRecorder()

	server.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("got status %d, want %d", rec.Code, http.StatusInternalServerError)
	}
	if !strings.Contains(rec.Body.String(), wantErr.Error()) {
		t.Fatalf("body %q should contain %q", rec.Body.String(), wantErr.Error())
	}
}

func TestServer_CustomErrorEncoder(t *testing.T) {
	ep := harmoniendpoint.Endpoint[int, int](func(ctx context.Context, req int) (int, error) {
		return 0, errors.New("not found")
	})

	dec := func(ctx context.Context, r *http.Request) (int, error) {
		return 1, nil
	}

	enc := func(ctx context.Context, w http.ResponseWriter, resp int) error {
		return json.NewEncoder(w).Encode(resp)
	}

	customErrEnc := func(ctx context.Context, w http.ResponseWriter, err error) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error":"not found"}`))
	}

	server := httptransport.NewServer(ep, dec, enc, httptransport.WithErrorEncoder[int, int](customErrEnc))

	req := httptest.NewRequest("POST", "/", strings.NewReader("1"))
	rec := httptest.NewRecorder()

	server.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("got status %d, want %d", rec.Code, http.StatusNotFound)
	}
}

func TestServer_BeforeHook(t *testing.T) {
	ep := harmoniendpoint.Endpoint[int, int](func(ctx context.Context, req int) (int, error) {
		v := ctx.Value("before-key")
		if v != "before-value" {
			t.Fatal("before hook did not set context value")
		}
		return req, nil
	})

	dec := func(ctx context.Context, r *http.Request) (int, error) {
		return 1, nil
	}

	enc := func(ctx context.Context, w http.ResponseWriter, resp int) error {
		return json.NewEncoder(w).Encode(resp)
	}

	beforeHook := func(ctx context.Context, r *http.Request) context.Context {
		return context.WithValue(ctx, "before-key", "before-value")
	}

	server := httptransport.NewServer(ep, dec, enc, httptransport.WithBefore[int, int](beforeHook))

	req := httptest.NewRequest("POST", "/", strings.NewReader("1"))
	rec := httptest.NewRecorder()

	server.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("got status %d, want %d", rec.Code, http.StatusOK)
	}
}

func TestServer_Shutdown(t *testing.T) {
	ep := harmoniendpoint.Endpoint[int, int](func(ctx context.Context, req int) (int, error) {
		return req, nil
	})

	dec := func(ctx context.Context, r *http.Request) (int, error) { return 1, nil }
	enc := func(ctx context.Context, w http.ResponseWriter, resp int) error { return nil }
	server := httptransport.NewServer(ep, dec, enc)

	// Shutdown on server that hasn't started should be nil.
	if err := server.Shutdown(context.Background()); err != nil {
		t.Fatalf("unexpected shutdown error: %v", err)
	}
}
