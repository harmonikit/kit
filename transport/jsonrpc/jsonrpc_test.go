package jsonrpc_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/harmonikit/harmoni/endpoint"
	jsonrpc "github.com/harmonikit/kit/transport/jsonrpc"
)

type addReq struct {
	A int `json:"a"`
	B int `json:"b"`
}

func TestServer_Success(t *testing.T) {
	ep := endpoint.Endpoint[addReq, int](func(ctx context.Context, req addReq) (int, error) {
		return req.A + req.B, nil
	})

	server := jsonrpc.NewServer(ep, "add")

	body := `{"jsonrpc":"2.0","method":"add","params":{"a":21,"b":21},"id":1}`
	req := httptest.NewRequest("POST", "/", strings.NewReader(body))
	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("got status %d, want 200", rec.Code)
	}

	var resp struct {
		Result int `json:"result"`
		ID     int `json:"id"`
	}
	json.NewDecoder(rec.Body).Decode(&resp)
	if resp.Result != 42 {
		t.Fatalf("got result %d, want 42", resp.Result)
	}
	if resp.ID != 1 {
		t.Fatalf("got id %d, want 1", resp.ID)
	}
}

func TestServer_InvalidJSON(t *testing.T) {
	ep := endpoint.Endpoint[int, int](func(ctx context.Context, req int) (int, error) {
		return req, nil
	})
	server := jsonrpc.NewServer(ep, "test")

	req := httptest.NewRequest("POST", "/", strings.NewReader("not json"))
	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)

	var resp struct {
		Error struct{ Code int } `json:"error"`
	}
	json.NewDecoder(rec.Body).Decode(&resp)
	if resp.Error.Code != jsonrpc.ErrParse {
		t.Fatalf("got code %d, want %d", resp.Error.Code, jsonrpc.ErrParse)
	}
}

func TestServer_WrongVersion(t *testing.T) {
	ep := endpoint.Endpoint[int, int](func(ctx context.Context, req int) (int, error) {
		return req, nil
	})
	server := jsonrpc.NewServer(ep, "test")

	body := `{"jsonrpc":"1.0","method":"test","id":1}`
	req := httptest.NewRequest("POST", "/", strings.NewReader(body))
	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)

	var resp struct {
		Error struct{ Code int } `json:"error"`
	}
	json.NewDecoder(rec.Body).Decode(&resp)
	if resp.Error.Code != jsonrpc.ErrInvalid {
		t.Fatalf("got code %d, want %d", resp.Error.Code, jsonrpc.ErrInvalid)
	}
}

func TestServer_MethodNotFound(t *testing.T) {
	ep := endpoint.Endpoint[int, int](func(ctx context.Context, req int) (int, error) {
		return req, nil
	})
	server := jsonrpc.NewServer(ep, "add")

	body := `{"jsonrpc":"2.0","method":"subtract","id":1}`
	req := httptest.NewRequest("POST", "/", strings.NewReader(body))
	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)

	var resp struct {
		Error struct{ Code int } `json:"error"`
	}
	json.NewDecoder(rec.Body).Decode(&resp)
	if resp.Error.Code != jsonrpc.ErrMethod {
		t.Fatalf("got code %d, want %d", resp.Error.Code, jsonrpc.ErrMethod)
	}
}

func TestServer_EndpointError(t *testing.T) {
	ep := endpoint.Endpoint[int, int](func(ctx context.Context, req int) (int, error) {
		return 0, fmt.Errorf("business error")
	})
	server := jsonrpc.NewServer(ep, "test")

	body := `{"jsonrpc":"2.0","method":"test","id":1}`
	req := httptest.NewRequest("POST", "/", strings.NewReader(body))
	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)

	var resp struct {
		Error struct{ Code int } `json:"error"`
	}
	json.NewDecoder(rec.Body).Decode(&resp)
	if resp.Error.Code != jsonrpc.ErrInternal {
		t.Fatalf("got code %d, want %d", resp.Error.Code, jsonrpc.ErrInternal)
	}
}

func TestServer_Notification(t *testing.T) {
	var called bool
	ep := endpoint.Endpoint[int, int](func(ctx context.Context, req int) (int, error) {
		called = true
		return req, nil
	})
	server := jsonrpc.NewServer(ep, "test")

	// Notification: no "id" field.
	body := `{"jsonrpc":"2.0","method":"test","params":42}`
	req := httptest.NewRequest("POST", "/", strings.NewReader(body))
	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)

	if !called {
		t.Fatal("endpoint should have been called for notification")
	}
	// No body for notification.
	if rec.Body.Len() > 0 {
		t.Fatal("notification should have no response body")
	}
}

func TestServer_ArrayParams(t *testing.T) {
	ep := endpoint.Endpoint[struct{ A, B int }, int](func(ctx context.Context, req struct{ A, B int }) (int, error) {
		return req.A + req.B, nil
	})
	server := jsonrpc.NewServer(ep, "add")

	body := `{"jsonrpc":"2.0","method":"add","params":[21,21],"id":1}`
	req := httptest.NewRequest("POST", "/", strings.NewReader(body))
	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("got status %d, want 200", rec.Code)
	}
}

func TestClient(t *testing.T) {
	// Start a test server that echoes the sum.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Params struct{ A, B int } `json:"params"`
			ID     int                `json:"id"`
		}
		json.NewDecoder(r.Body).Decode(&req)
		resp := map[string]any{
			"jsonrpc": "2.0",
			"result":  req.Params.A + req.Params.B,
			"id":      req.ID,
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer ts.Close()

	client := jsonrpc.NewClient[struct{ A, B int }, int](ts.URL, "add")
	ep := client.Endpoint()

	resp, err := ep(context.Background(), struct{ A, B int }{21, 21})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp != 42 {
		t.Fatalf("got %d, want 42", resp)
	}
}
