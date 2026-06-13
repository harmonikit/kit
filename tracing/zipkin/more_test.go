package zipkin_test

import (
	"context"
	"testing"

	harmonitracing "github.com/harmonikit/harmoni/tracing"
	zipkintrace "github.com/harmonikit/kit/tracing/zipkin"
)

func TestTracer_End_WrongSpanType(t *testing.T) {
	tracer := zipkintrace.NewTracer[int, string]()
	ctx := context.Background()
	// Pass a non-Zipkin span — should be a no-op.
	tracer.End(ctx, harmonitracing.NopSpan{}, "ok", nil)
}

func TestSpan_SetAttributes_Direct(t *testing.T) {
	s := &zipkintrace.Span{}
	s.SetAttributes("key", "value")
	s.End()
}
