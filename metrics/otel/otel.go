// Package otel adapts the OpenTelemetry metrics SDK to the harmoni/metrics
// interfaces.
//
// These are stub implementations. A production version requires an OTEL
// MeterProvider to create real instruments.
package otel

import (
	"sync/atomic"

	harmonimetric "github.com/harmonikit/harmoni/metrics"
)

// Compile-time interface conformance.
var (
	_ harmonimetric.Counter   = (*Counter)(nil)
	_ harmonimetric.Gauge     = (*Gauge)(nil)
	_ harmonimetric.Histogram = (*Histogram)(nil)
)

// Counter implements harmonimetric.Counter backed by an atomic float64.
type Counter struct {
	val atomic.Int64
	lvs []string
}

// NewCounter returns a Counter.
func NewCounter() *Counter {
	return &Counter{val: atomic.Int64{}, lvs: nil}
}

// With returns a new Counter with the given label values.
func (c *Counter) With(labelValues ...string) harmonimetric.Counter {
	return &Counter{val: atomic.Int64{}, lvs: append(c.lvs, labelValues...)}
}

// Add increments the counter by delta.
func (c *Counter) Add(delta float64) {
	c.val.Add(int64(delta))
}

// Value returns the current counter value (for testing).
func (c *Counter) Value() float64 {
	return float64(c.val.Load())
}

// Gauge implements harmonimetric.Gauge backed by an atomic float64.
type Gauge struct {
	val atomic.Int64
	lvs []string
}

// NewGauge returns a Gauge.
func NewGauge() *Gauge {
	return &Gauge{val: atomic.Int64{}, lvs: nil}
}

// With returns a new Gauge with the given label values.
func (g *Gauge) With(labelValues ...string) harmonimetric.Gauge {
	return &Gauge{val: atomic.Int64{}, lvs: append(g.lvs, labelValues...)}
}

// Set sets the gauge to an absolute value.
func (g *Gauge) Set(value float64) {
	g.val.Store(int64(value))
}

// Add increments or decrements the gauge by delta.
func (g *Gauge) Add(delta float64) {
	g.val.Add(int64(delta))
}

// Value returns the current gauge value (for testing).
func (g *Gauge) Value() float64 {
	return float64(g.val.Load())
}

// Histogram implements harmonimetric.Histogram backed by a simple value store.
type Histogram struct {
	vals []float64
	lvs  []string
}

// NewHistogram returns a Histogram.
func NewHistogram() *Histogram {
	return &Histogram{vals: nil, lvs: nil}
}

// With returns a new Histogram with the given label values.
func (h *Histogram) With(labelValues ...string) harmonimetric.Histogram {
	return &Histogram{vals: nil, lvs: append(h.lvs, labelValues...)}
}

// Observe records a value in the histogram.
func (h *Histogram) Observe(value float64) {
	h.vals = append(h.vals, value)
}

// Observations returns recorded values (for testing).
func (h *Histogram) Observations() []float64 {
	return h.vals
}
