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
// This is a stub — production code should use an OTEL Float64Counter.
type Counter struct {
	val  atomic.Int64
	lvs  []string
}

// NewCounter returns a Counter.
func NewCounter() *Counter {
	return &Counter{}
}

func (c *Counter) With(labelValues ...string) harmonimetric.Counter {
	return &Counter{lvs: append(c.lvs, labelValues...)}
}

func (c *Counter) Add(delta float64) {
	c.val.Add(int64(delta))
}

// Value returns the current counter value (for testing).
func (c *Counter) Value() float64 {
	return float64(c.val.Load())
}

// Gauge implements harmonimetric.Gauge backed by an atomic float64.
type Gauge struct {
	val  atomic.Int64
	lvs  []string
}

// NewGauge returns a Gauge.
func NewGauge() *Gauge {
	return &Gauge{}
}

func (g *Gauge) With(labelValues ...string) harmonimetric.Gauge {
	return &Gauge{lvs: append(g.lvs, labelValues...)}
}

func (g *Gauge) Set(value float64) {
	g.val.Store(int64(value))
}

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
	return &Histogram{}
}

func (h *Histogram) With(labelValues ...string) harmonimetric.Histogram {
	return &Histogram{lvs: append(h.lvs, labelValues...)}
}

func (h *Histogram) Observe(value float64) {
	h.vals = append(h.vals, value)
}

// Observations returns recorded values (for testing).
func (h *Histogram) Observations() []float64 {
	return h.vals
}