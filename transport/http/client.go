package http

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"

	harmoniendpoint "github.com/harmonikit/harmoni/endpoint"
)

// Client is an HTTP client that calls a remote endpoint.
type Client[Req, Resp any] struct {
	client  *http.Client
	method  string
	url     *url.URL
	enc     EncodeRequestFunc[Req]
	dec     DecodeResponseFunc[Resp]
}

// EncodeRequestFunc encodes a domain request into an HTTP request body.
type EncodeRequestFunc[Req any] func(ctx context.Context, req Req) (io.Reader, error)

// DecodeResponseFunc decodes an HTTP response body into a domain response.
type DecodeResponseFunc[Resp any] func(ctx context.Context, r *http.Response) (Resp, error)

// ClientOption configures a Client.
type ClientOption[Req, Resp any] func(*Client[Req, Resp])

// WithHTTPClient sets a custom http.Client.
func WithHTTPClient[Req, Resp any](c *http.Client) ClientOption[Req, Resp] {
	return func(cl *Client[Req, Resp]) { cl.client = c }
}

// NewClient creates an HTTP client for calling a remote service.
func NewClient[Req, Resp any](
	method string,
	target *url.URL,
	enc EncodeRequestFunc[Req],
	dec DecodeResponseFunc[Resp],
	opts ...ClientOption[Req, Resp],
) *Client[Req, Resp] {
	c := &Client[Req, Resp]{
		client: http.DefaultClient,
		method: method,
		url:    target,
		enc:    enc,
		dec:    dec,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// Endpoint returns an endpoint.Endpoint that calls the remote HTTP service.
func (c *Client[Req, Resp]) Endpoint() harmoniendpoint.Endpoint[Req, Resp] {
	return func(ctx context.Context, req Req) (Resp, error) {
		body, err := c.enc(ctx, req)
		if err != nil {
			var zero Resp
			return zero, fmt.Errorf("encode request: %w", err)
		}

		// Determine the body reader.
		var r io.Reader
		if body != nil {
			r = body
		}

		httpReq, err := http.NewRequestWithContext(ctx, c.method, c.url.String(), r)
		if err != nil {
			var zero Resp
			return zero, fmt.Errorf("create request: %w", err)
		}

		httpResp, err := c.client.Do(httpReq)
		if err != nil {
			var zero Resp
			return zero, fmt.Errorf("do request: %w", err)
		}
		defer httpResp.Body.Close()

		if httpResp.StatusCode < 200 || httpResp.StatusCode >= 300 {
			// Read body for error context.
			bodyBytes, _ := io.ReadAll(httpResp.Body)
			var zero Resp
			return zero, fmt.Errorf("http %d: %s", httpResp.StatusCode, string(bytes.TrimSpace(bodyBytes)))
		}

		resp, err := c.dec(ctx, httpResp)
		if err != nil {
			var zero Resp
			return zero, fmt.Errorf("decode response: %w", err)
		}
		return resp, nil
	}
}
