package opentelemetry_test

import (
	"context"
	"testing"

	harmonitracing "github.com/harmonikit/harmoni/tracing"
	oteltrace "github.com/harmonikit/kit/tracing/opentelemetry"
)

func TestTracer_End_WrongSpanType(t *testing.T) {
	tracer := oteltrace.NewTracer[int, string]("test")

	ctx := context.Background()
	// Pass a non-OTEL span — should be a no-op, not a panic.
	tracer.End(ctx, harmonitracing.NopSpan{}, "ok", nil)
}
