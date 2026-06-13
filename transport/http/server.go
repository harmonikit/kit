package http

import (
	"context"
	"net/http"

	"github.com/harmonikit/harmoni/endpoint"
)

// Server is an HTTP transport server. It wraps an endpoint and serves it over HTTP.
type Server[Req, Resp any] struct {
	server   *http.Server
	endpoint endpoint.Endpoint[Req, Resp]
	dec      DecodeRequestFunc[Req]
	enc      EncodeResponseFunc[Resp]
	encError EncodeErrorFunc
	before   []BeforeFunc
}

// BeforeFunc is a hook executed before each request. If it returns a non-nil
// context, that context is used for the endpoint call. If it writes a response
// and returns nil context, the endpoint is not called.
type BeforeFunc func(ctx context.Context, r *http.Request) context.Context

// ServerOption configures a Server.
type ServerOption[Req, Resp any] func(*Server[Req, Resp])

// WithErrorEncoder sets a custom error encoder.
func WithErrorEncoder[Req, Resp any](fn EncodeErrorFunc) ServerOption[Req, Resp] {
	return func(s *Server[Req, Resp]) { s.encError = fn }
}

// WithBefore adds a before-request hook.
func WithBefore[Req, Resp any](fn BeforeFunc) ServerOption[Req, Resp] {
	return func(s *Server[Req, Resp]) { s.before = append(s.before, fn) }
}

// defaultErrorEncoder writes a 500 status with the error message.
// For status codes that must not have a body (1xx, 204, 304), the body
// is omitted as required by RFC 7230.
func defaultErrorEncoder(_ context.Context, w http.ResponseWriter, err error) {
	code := http.StatusInternalServerError
	if !allowsBody(code) {
		w.WriteHeader(code)
		return
	}
	http.Error(w, err.Error(), code)
}

// allowsBody reports whether the HTTP status code permits a response body.
// 1xx informational codes, 204 No Content, and 304 Not Modified must not
// include a body per RFC 7230.
func allowsBody(code int) bool {
	return code >= 200 && code != http.StatusNoContent && code != http.StatusNotModified
}

// NewServer creates an HTTP handler from an endpoint, request decoder, and
// response encoder.
func NewServer[Req, Resp any](
	ep endpoint.Endpoint[Req, Resp],
	dec DecodeRequestFunc[Req],
	enc EncodeResponseFunc[Resp],
	opts ...ServerOption[Req, Resp],
) *Server[Req, Resp] {
	s := &Server[Req, Resp]{
		server:   nil,
		endpoint: ep,
		dec:      dec,
		enc:      enc,
		encError: defaultErrorEncoder,
		before:   nil,
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

// ServeHTTP implements http.Handler.
func (s *Server[Req, Resp]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Wrap to prevent body on no-body status codes (RFC 7230).
	rw := &recordsStatus{ResponseWriter: w}

	// Run before hooks.
	for _, fn := range s.before {
		ctx = fn(ctx, r)
	}

	req, err := s.dec(ctx, r)
	if err != nil {
		s.encError(ctx, rw, err)
		return
	}

	resp, err := s.endpoint(ctx, req)
	if err != nil {
		s.encError(ctx, rw, err)
		return
	}

	if err := s.enc(ctx, rw, resp); err != nil {
		s.encError(ctx, rw, err)
		return
	}
}

// recordsStatus wraps http.ResponseWriter to capture the status code.
type recordsStatus struct {
	http.ResponseWriter
	wroteHeader bool
	status      int
}

func (r *recordsStatus) WriteHeader(code int) {
	if !r.wroteHeader {
		r.status = code
		r.wroteHeader = true
		if !allowsBody(code) {
			r.ResponseWriter.WriteHeader(code)
			return
		}
	}
	r.ResponseWriter.WriteHeader(code)
}

func (r *recordsStatus) Write(b []byte) (int, error) {
	if !r.wroteHeader {
		r.WriteHeader(http.StatusOK)
	}
	if !allowsBody(r.status) {
		return len(b), nil // discard body for no-body status codes
	}
	return r.ResponseWriter.Write(b)
}

// ListenAndServe starts the HTTP server on the given address.
func (s *Server[Req, Resp]) ListenAndServe(addr string) error {
	//nolint:exhaustruct
	s.server = &http.Server{
		Addr:    addr,
		Handler: s,
	}
	return s.server.ListenAndServe()
}

// Shutdown gracefully shuts down the server.
func (s *Server[Req, Resp]) Shutdown(ctx context.Context) error {
	if s.server == nil {
		return nil
	}
	return s.server.Shutdown(ctx)
}
