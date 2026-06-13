// Package opentelemetry adapts go.opentelemetry.io/otel to the harmoni/tracing
// interfaces.
package opentelemetry

import (
	"context"
	"fmt"

	harmonitracing "github.com/harmonikit/harmoni/tracing"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// Tracer implements harmonitracing.Tracer using OpenTelemetry.
type Tracer[Req, Resp any] struct {
	tracer trace.Tracer
}

// NewTracer returns a Tracer backed by an OTEL tracer.
// name is typically the service or instrumentation name.
func NewTracer[Req, Resp any](name string, opts ...trace.TracerOption) *Tracer[Req, Resp] {
	return &Tracer[Req, Resp]{
		tracer: otel.Tracer(name, opts...),
	}
}

// Start begins a new span and returns the updated context.
func (t *Tracer[Req, Resp]) Start(ctx context.Context, operationName string, req Req) (context.Context, harmonitracing.Span) {
	ctx, sp := t.tracer.Start(ctx, operationName)
	return ctx, &Span{span: sp}
}

// End completes the span, recording the response and error.
func (t *Tracer[Req, Resp]) End(ctx context.Context, span harmonitracing.Span, resp Resp, err error) {
	sp, ok := span.(*Span)
	if !ok {
		return
	}
	if err != nil {
		sp.span.SetStatus(codes.Error, err.Error())
		sp.span.RecordError(err)
	}
	sp.span.SetAttributes(attribute.String("response", fmt.Sprintf("%v", resp)))
	sp.span.End()
}

// Span wraps an OTEL span.
type Span struct {
	span trace.Span
}

func (s *Span) End() {
	s.span.End()
}

func (s *Span) SetAttributes(attrs ...any) {
	otelAttrs := make([]attribute.KeyValue, 0, len(attrs)/2)
	for i := 0; i < len(attrs); i += 2 {
		if i+1 < len(attrs) {
			otelAttrs = append(otelAttrs, attribute.String(fmt.Sprintf("%v", attrs[i]), fmt.Sprintf("%v", attrs[i+1])))
		}
	}
	s.span.SetAttributes(otelAttrs...)
}
