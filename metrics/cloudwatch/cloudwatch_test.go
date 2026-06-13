package cloudwatch_test

import (
	"testing"

	harmonimetric "github.com/harmonikit/harmoni/metrics"
	cloudwatch "github.com/harmonikit/kit/metrics/cloudwatch"
)

func TestCounter(t *testing.T) {
	c := cloudwatch.NewCounter("MyApp", "requests", dim("Environment", "prod"))
	c.Add(1)
	c.Add(5)
	if c.Value() != 6 {
		t.Fatalf("got %f, want 6", c.Value())
	}
}

func TestCounter_With(t *testing.T) {
	c := cloudwatch.NewCounter("MyApp", "requests")
	labeled := c.With("method", "GET")
	labeled.Add(1)
	if c.Value() != 0 {
		t.Fatal("With should not mutate original")
	}
}

func TestGauge(t *testing.T) {
	g := cloudwatch.NewGauge("MyApp", "memory")
	g.Set(100)
	g.Add(-20)
	if g.Value() != 80 {
		t.Fatalf("got %f, want 80", g.Value())
	}
}

func TestHistogram(t *testing.T) {
	h := cloudwatch.NewHistogram("MyApp", "latency")
	h.Observe(1.0)
	h.Observe(2.0)
	h.Observe(3.0)
	if len(h.Observations()) != 3 {
		t.Fatalf("got %d, want 3", len(h.Observations()))
	}
}

func TestInterfaces(t *testing.T) {
	var _ harmonimetric.Counter = cloudwatch.NewCounter("ns", "n")
	var _ harmonimetric.Gauge = cloudwatch.NewGauge("ns", "n")
	var _ harmonimetric.Histogram = cloudwatch.NewHistogram("ns", "n")
}

func dim(name, value string) cloudwatch.Dimension {
	return cloudwatch.Dimension{Name: name, Value: value}
}
