package dogstatsd_test

import (
	"testing"

	harmonimetric "github.com/harmonikit/harmoni/metrics"
	dogstatsd "github.com/harmonikit/kit/metrics/dogstatsd"
)

func TestCounter(t *testing.T) {
	c := dogstatsd.NewCounter("requests", "service:api")
	c.Add(1)
	c.Add(5)
	if c.Value() != 6 {
		t.Fatalf("got %f, want 6", c.Value())
	}
}

func TestCounter_With(t *testing.T) {
	c := dogstatsd.NewCounter("requests")
	labeled := c.With("method", "GET")
	labeled.Add(1)
	if c.Value() != 0 {
		t.Fatal("With should not mutate original")
	}
}

func TestGauge(t *testing.T) {
	g := dogstatsd.NewGauge("memory")
	g.Set(100)
	g.Add(-20)
	if g.Value() != 80 {
		t.Fatalf("got %f, want 80", g.Value())
	}
}

func TestHistogram(t *testing.T) {
	h := dogstatsd.NewHistogram("latency")
	h.Observe(1.0)
	h.Observe(2.0)
	h.Observe(3.0)
	if len(h.Observations()) != 3 {
		t.Fatalf("got %d, want 3", len(h.Observations()))
	}
}

func TestInterfaces(t *testing.T) {
	var _ harmonimetric.Counter = dogstatsd.NewCounter("test")
	var _ harmonimetric.Gauge = dogstatsd.NewGauge("test")
	var _ harmonimetric.Histogram = dogstatsd.NewHistogram("test")
}
