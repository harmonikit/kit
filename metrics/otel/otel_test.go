package otel_test

import (
	"testing"

	harmonimetric "github.com/harmonikit/harmoni/metrics"
	otelmetrics "github.com/harmonikit/kit/metrics/otel"
)

func TestCounter_Add(t *testing.T) {
	c := otelmetrics.NewCounter()
	c.Add(1)
	c.Add(5)
	if c.Value() != 6 {
		t.Fatalf("got %f, want 6", c.Value())
	}
}

func TestCounter_With(t *testing.T) {
	c := otelmetrics.NewCounter()
	labeled := c.With("method", "GET")
	labeled.Add(1)
	if c.Value() != 0 {
		t.Fatal("With should create a new counter, not mutate the original")
	}
}

func TestGauge_Set(t *testing.T) {
	g := otelmetrics.NewGauge()
	g.Set(42)
	if g.Value() != 42 {
		t.Fatalf("got %f, want 42", g.Value())
	}
}

func TestGauge_Add(t *testing.T) {
	g := otelmetrics.NewGauge()
	g.Add(10)
	g.Add(-3)
	if g.Value() != 7 {
		t.Fatalf("got %f, want 7", g.Value())
	}
}

func TestHistogram_Observe(t *testing.T) {
	h := otelmetrics.NewHistogram()
	h.Observe(1.0)
	h.Observe(2.0)
	h.Observe(3.0)
	if len(h.Observations()) != 3 {
		t.Fatalf("got %d observations, want 3", len(h.Observations()))
	}
}

func TestHistogram_With(t *testing.T) {
	h := otelmetrics.NewHistogram()
	labeled := h.With("path", "/api")
	labeled.Observe(1.0)
	if len(h.Observations()) != 0 {
		t.Fatal("With should create a new histogram")
	}
}

func TestInterfaces(t *testing.T) {
	var _ harmonimetric.Counter = otelmetrics.NewCounter()
	var _ harmonimetric.Gauge = otelmetrics.NewGauge()
	var _ harmonimetric.Histogram = otelmetrics.NewHistogram()
}