// Package cloudwatch adapts AWS CloudWatch to the harmoni/metrics interfaces.
//
// This is a stub implementation that tracks values in memory for testing.
// A production version depends on github.com/aws/aws-sdk-go-v2/service/cloudwatch
// and emits PutMetricData calls.
package cloudwatch

import (
	"sync"

	harmonimetric "github.com/harmonikit/harmoni/metrics"
)

// Compile-time interface conformance.
var (
	_ harmonimetric.Counter   = (*Counter)(nil)
	_ harmonimetric.Gauge     = (*Gauge)(nil)
	_ harmonimetric.Histogram = (*Histogram)(nil)
)

// Dimension is a CloudWatch metric dimension.
type Dimension struct{ Name, Value string }

// Counter implements harmonimetric.Counter for CloudWatch.
type Counter struct {
	mu         sync.Mutex
	namespace  string
	name       string
	dimensions []Dimension
	value      float64
}

// NewCounter returns a Counter backed by CloudWatch metrics.
func NewCounter(namespace, name string, dims ...Dimension) *Counter {
	return &Counter{namespace: namespace, name: name, dimensions: dims}
}

// With returns a new Counter with additional dimensions.
func (c *Counter) With(labelValues ...string) harmonimetric.Counter {
	dims := make([]Dimension, len(c.dimensions))
	copy(dims, c.dimensions)
	for i := 0; i+1 < len(labelValues); i += 2 {
		dims = append(dims, Dimension{Name: labelValues[i], Value: labelValues[i+1]})
	}
	return &Counter{namespace: c.namespace, name: c.name, dimensions: dims}
}

// Add increments the counter by delta.
func (c *Counter) Add(delta float64) {
	c.mu.Lock()
	c.value += delta
	c.mu.Unlock()
}

// Value returns the current counter value (for testing).
func (c *Counter) Value() float64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.value
}

// Gauge implements harmonimetric.Gauge for CloudWatch.
type Gauge struct {
	mu         sync.Mutex
	namespace  string
	name       string
	dimensions []Dimension
	value      float64
}

// NewGauge returns a Gauge backed by CloudWatch metrics.
func NewGauge(namespace, name string, dims ...Dimension) *Gauge {
	return &Gauge{namespace: namespace, name: name, dimensions: dims}
}

// With returns a new Gauge with additional dimensions.
func (g *Gauge) With(labelValues ...string) harmonimetric.Gauge {
	dims := make([]Dimension, len(g.dimensions))
	copy(dims, g.dimensions)
	for i := 0; i+1 < len(labelValues); i += 2 {
		dims = append(dims, Dimension{Name: labelValues[i], Value: labelValues[i+1]})
	}
	return &Gauge{namespace: g.namespace, name: g.name, dimensions: dims}
}

// Set sets the gauge to an absolute value.
func (g *Gauge) Set(value float64) {
	g.mu.Lock()
	g.value = value
	g.mu.Unlock()
}

// Add increments or decrements the gauge by delta.
func (g *Gauge) Add(delta float64) {
	g.mu.Lock()
	g.value += delta
	g.mu.Unlock()
}

// Value returns the current gauge value (for testing).
func (g *Gauge) Value() float64 {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.value
}

// Histogram implements harmonimetric.Histogram for CloudWatch.
type Histogram struct {
	mu         sync.Mutex
	namespace  string
	name       string
	dimensions []Dimension
	vals       []float64
}

// NewHistogram returns a Histogram backed by CloudWatch metrics.
func NewHistogram(namespace, name string, dims ...Dimension) *Histogram {
	return &Histogram{namespace: namespace, name: name, dimensions: dims}
}

// With returns a new Histogram with additional dimensions.
func (h *Histogram) With(labelValues ...string) harmonimetric.Histogram {
	dims := make([]Dimension, len(h.dimensions))
	copy(dims, h.dimensions)
	for i := 0; i+1 < len(labelValues); i += 2 {
		dims = append(dims, Dimension{Name: labelValues[i], Value: labelValues[i+1]})
	}
	return &Histogram{namespace: h.namespace, name: h.name, dimensions: dims}
}

// Observe records a value in the histogram.
func (h *Histogram) Observe(value float64) {
	h.mu.Lock()
	h.vals = append(h.vals, value)
	h.mu.Unlock()
}

// Observations returns recorded values (for testing).
func (h *Histogram) Observations() []float64 {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.vals
}
