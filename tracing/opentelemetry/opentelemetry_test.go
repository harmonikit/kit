package opentelemetry_test

import (
	"context"
	"errors"
	"testing"

	harmonitracing "github.com/harmonikit/harmoni/tracing"
	oteltrace "github.com/harmonikit/kit/tracing/opentelemetry"
)

func TestTracer_StartEnd(t *testing.T) {
	tracer := oteltrace.NewTracer[int, string]("test")

	ctx := context.Background()
	newCtx, span := tracer.Start(ctx, "operation", 42)
	if newCtx == nil {
		t.Fatal("expected non-nil context")
	}
	if span == nil {
		t.Fatal("expected non-nil span")
	}

	tracer.End(newCtx, span, "ok", nil)
}

func TestTracer_End_WithError(t *testing.T) {
	tracer := oteltrace.NewTracer[int, string]("test")

	ctx := context.Background()
	newCtx, span := tracer.Start(ctx, "operation", 42)

	err := errors.New("something failed")
	tracer.End(newCtx, span, "partial", err)
}

func TestSpan_SetAttributes(t *testing.T) {
	tracer := oteltrace.NewTracer[int, string]("test")

	_, span := tracer.Start(context.Background(), "op", 1)
	span.SetAttributes("key", "value", "count", 42)
	span.End()
}

func TestTracer_Interface(t *testing.T) {
	var _ harmonitracing.Tracer[int, string] = oteltrace.NewTracer[int, string]("test")
}