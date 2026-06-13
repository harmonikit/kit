package zipkin_test

import (
	"context"
	"errors"
	"testing"

	harmonitracing "github.com/harmonikit/harmoni/tracing"
	zipkintrace "github.com/harmonikit/kit/tracing/zipkin"
)

func TestTracer_StartEnd(t *testing.T) {
	tracer := zipkintrace.NewTracer[int, string]()

	ctx := context.Background()
	newCtx, span := tracer.Start(ctx, "operation", 42)
	if newCtx == nil {
		t.Fatal("expected non-nil context")
	}
	sp := span.(*zipkintrace.Span)
	if sp == nil {
		t.Fatal("expected non-nil span")
	}

	tracer.End(newCtx, span, "ok", nil)
}

func TestTracer_End_WithError(t *testing.T) {
	tracer := zipkintrace.NewTracer[int, string]()

	ctx := context.Background()
	_, span := tracer.Start(ctx, "operation", 42)
	err := errors.New("fail")
	tracer.End(ctx, span, "", err)
}

func TestSpan_SetAttributes(t *testing.T) {
	tracer := zipkintrace.NewTracer[int, string]()

	_, span := tracer.Start(context.Background(), "op", 1)
	span.SetAttributes("k1", "v1", "k2", 42)
	span.End()
}

func TestTracer_Interface(t *testing.T) {
	var _ harmonitracing.Tracer[int, string] = zipkintrace.NewTracer[int, string]()
}
