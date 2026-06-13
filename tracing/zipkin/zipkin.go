// Package zipkin adapts Zipkin tracing to the harmoni/tracing interfaces.
//
// This is a stub implementation. A production version would use
// go.opentelemetry.io/otel/exporters/zipkin or the legacy
// github.com/openzipkin/zipkin-go reporter.
package zipkin

import (
	"context"
	"fmt"

	harmonitracing "github.com/harmonikit/harmoni/tracing"
)

// Tracer implements harmonitracing.Tracer for Zipkin spans.
type Tracer[Req, Resp any] struct{}

// NewTracer returns a Zipkin-backed Tracer.
func NewTracer[Req, Resp any]() *Tracer[Req, Resp] {
	return &Tracer[Req, Resp]{}
}

// Start begins a new span.
func (t *Tracer[Req, Resp]) Start(ctx context.Context, operationName string, req Req) (context.Context, harmonitracing.Span) {
	return ctx, &Span{name: operationName}
}

// End completes the span.
func (t *Tracer[Req, Resp]) End(ctx context.Context, span harmonitracing.Span, resp Resp, err error) {
	sp, ok := span.(*Span)
	if !ok {
		return
	}
	sp.finished = true
	if err != nil {
		sp.err = err
	}
}

// Span represents a Zipkin span.
type Span struct {
	name     string
	finished bool
	err      error
	attrs    []string
}

func (s *Span) End() {
	s.finished = true
}

func (s *Span) SetAttributes(attrs ...any) {
	for _, a := range attrs {
		s.attrs = append(s.attrs, fmt.Sprintf("%v", a))
	}
}
