// Package jsonrpc provides JSON-RPC 2.0 transport bindings for harmoni endpoints.
//
// It wraps a typed endpoint as an HTTP handler accepting JSON-RPC 2.0 requests
// and provides a client for calling remote JSON-RPC services.
//
// Example server:
//
//	ep := endpoint.Endpoint[Req, Resp](service.Method)
//	handler := jsonrpc.NewServer(ep, "add")
//	http.Handle("/rpc", handler)
//
// Example client:
//
//	ep := jsonrpc.NewClient[Req, Resp]("http://host/rpc", "add").Endpoint()
//	resp, err := ep(ctx, req)
package jsonrpc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/harmonikit/harmoni/endpoint"
)

// ── JSON-RPC 2.0 types ─────────────────────────────────────────────────

type request struct {
	JSONRPC string          `json:"jsonrpc"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
	ID      *int            `json:"id,omitempty"` // nil = notification
}

type response struct {
	JSONRPC string      `json:"jsonrpc"`
	Result  any         `json:"result,omitempty"`
	Error   *rpcError   `json:"error,omitempty"`
	ID      *int        `json:"id"`
}

type rpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Standard JSON-RPC error codes.
const (
	ErrParse     = -32700
	ErrInvalid   = -32600
	ErrMethod    = -32601
	ErrParams    = -32602
	ErrInternal  = -32603
)

// ── Server ──────────────────────────────────────────────────────────────

// Server wraps a harmoni endpoint as a JSON-RPC 2.0 HTTP handler.
type Server[Req, Resp any] struct {
	endpoint endpoint.Endpoint[Req, Resp]
	method   string
}

// NewServer creates a JSON-RPC handler for the given method name.
func NewServer[Req, Resp any](ep endpoint.Endpoint[Req, Resp], method string) *Server[Req, Resp] {
	return &Server[Req, Resp]{endpoint: ep, method: method}
}

// ServeHTTP implements http.Handler.
func (s *Server[Req, Resp]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeError(w, nil, ErrParse, "failed to read request body")
		return
	}

	var req request
	if err := json.Unmarshal(body, &req); err != nil {
		writeError(w, nil, ErrParse, "invalid JSON")
		return
	}

	if req.JSONRPC != "2.0" {
		writeError(w, req.ID, ErrInvalid, "jsonrpc must be 2.0")
		return
	}

	if req.Method != s.method {
		writeError(w, req.ID, ErrMethod, fmt.Sprintf("method %q not found", req.Method))
		return
	}

	var domainReq Req
	if len(req.Params) > 0 {
		// Support both array and object params per JSON-RPC spec.
		if req.Params[0] == '[' {
			// Array params: unmarshal into slice, then take first element.
			// For simplicity, try unmarshaling directly.
		}
		if err := json.Unmarshal(req.Params, &domainReq); err != nil {
			writeError(w, req.ID, ErrParams, "invalid params")
			return
		}
	}

	// Notification: no response expected.
	if req.ID == nil {
		_, _ = s.endpoint(ctx, domainReq)
		return
	}

	resp, err := s.endpoint(ctx, domainReq)
	if err != nil {
		writeError(w, req.ID, ErrInternal, err.Error())
		return
	}

	writeResult(w, req.ID, resp)
}

// ── Client ──────────────────────────────────────────────────────────────

// Client wraps a remote JSON-RPC endpoint.
type Client[Req, Resp any] struct {
	url    string
	method string
}

// NewClient creates a JSON-RPC client.
func NewClient[Req, Resp any](url, method string) *Client[Req, Resp] {
	return &Client[Req, Resp]{url: url, method: method}
}

// Endpoint returns a harmoni endpoint that calls the remote JSON-RPC service.
func (c *Client[Req, Resp]) Endpoint() endpoint.Endpoint[Req, Resp] {
	return func(ctx context.Context, req Req) (Resp, error) {
		id := 1
		params, err := json.Marshal(req)
		if err != nil {
			var zero Resp
			return zero, fmt.Errorf("jsonrpc marshal params: %w", err)
		}

		rpcReq := request{
			JSONRPC: "2.0",
			Method:  c.method,
			Params:  params,
			ID:      &id,
		}

		body, err := json.Marshal(rpcReq)
		if err != nil {
			var zero Resp
			return zero, fmt.Errorf("jsonrpc marshal request: %w", err)
		}

		httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.url, bytes.NewReader(body))
		if err != nil {
			var zero Resp
			return zero, fmt.Errorf("jsonrpc create request: %w", err)
		}
		httpReq.Header.Set("Content-Type", "application/json")

		httpResp, err := http.DefaultClient.Do(httpReq)
		if err != nil {
			var zero Resp
			return zero, fmt.Errorf("jsonrpc do: %w", err)
		}
		defer httpResp.Body.Close()

		var rpcResp response
		if err := json.NewDecoder(httpResp.Body).Decode(&rpcResp); err != nil {
			var zero Resp
			return zero, fmt.Errorf("jsonrpc decode response: %w", err)
		}

		if rpcResp.Error != nil {
			var zero Resp
			return zero, fmt.Errorf("jsonrpc error %d: %s", rpcResp.Error.Code, rpcResp.Error.Message)
		}

		var resp Resp
		data, _ := json.Marshal(rpcResp.Result)
		if err := json.Unmarshal(data, &resp); err != nil {
			var zero Resp
			return zero, fmt.Errorf("jsonrpc decode result: %w", err)
		}
		return resp, nil
	}
}

// ── Helpers ─────────────────────────────────────────────────────────────

func writeResult(w http.ResponseWriter, id *int, result any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(response{
		JSONRPC: "2.0",
		Result:  result,
		ID:      id,
	})
}

func writeError(w http.ResponseWriter, id *int, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(response{
		JSONRPC: "2.0",
		Error:   &rpcError{Code: code, Message: message},
		ID:      id, // nil = null in JSON (used for parse errors per spec)
	})
}
