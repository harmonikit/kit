package http_test

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	harmoniendpoint "github.com/harmonikit/harmoni/endpoint"
	httptransport "github.com/harmonikit/kit/transport/http"
)

func TestServer_NoBodyStatusCodes(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
	}{
		{"204 No Content", http.StatusNoContent},
		{"304 Not Modified", http.StatusNotModified},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ep := harmoniendpoint.Endpoint[int, int](func(ctx context.Context, req int) (int, error) {
				return 0, errors.New("no body error")
			})

			dec := func(ctx context.Context, r *http.Request) (int, error) { return 1, nil }
			enc := func(ctx context.Context, w http.ResponseWriter, resp int) error { return nil }

			noBodyEnc := func(ctx context.Context, w http.ResponseWriter, err error) {
				w.WriteHeader(tt.statusCode)
			}

			server := httptransport.NewServer(ep, dec, enc,
				httptransport.WithErrorEncoder[int, int](noBodyEnc))

			req := httptest.NewRequest("POST", "/", nil)
			rec := httptest.NewRecorder()
			server.ServeHTTP(rec, req)

			if rec.Code != tt.statusCode {
				t.Fatalf("got status %d, want %d", rec.Code, tt.statusCode)
			}
			body, _ := io.ReadAll(rec.Result().Body)
			if len(body) != 0 {
				t.Fatalf("expected empty body for %d, got: %q", tt.statusCode, body)
			}
		})
	}
}

func TestServer_AllowsBody_500(t *testing.T) {
	// 500 should have a body.
	ep := harmoniendpoint.Endpoint[int, int](func(ctx context.Context, req int) (int, error) {
		return 0, errors.New("internal error")
	})

	dec := func(ctx context.Context, r *http.Request) (int, error) { return 1, nil }
	enc := func(ctx context.Context, w http.ResponseWriter, resp int) error { return nil }

	server := httptransport.NewServer(ep, dec, enc)

	req := httptest.NewRequest("POST", "/", nil)
	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("got status %d, want %d", rec.Code, http.StatusInternalServerError)
	}
	body, _ := io.ReadAll(rec.Result().Body)
	if len(body) == 0 {
		t.Fatal("expected body for 500, got empty")
	}
}
