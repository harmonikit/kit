// Package dogstatsd adapts DataDog DogStatsD to the harmoni/metrics interfaces.
//
// This is a stub implementation that tracks values in memory for testing.
// A production version depends on github.com/DataDog/datadog-go/v5/statsd.
package dogstatsd

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

// Counter implements harmonimetric.Counter for DogStatsD.
type Counter struct {
	mu    sync.Mutex
	name  string
	tags  []string
	value float64
}

// NewCounter returns a Counter that emits DogStatsD count metrics.
func NewCounter(name string, tags ...string) *Counter {
	return &Counter{name: name, tags: tags}
}

// With returns a new Counter with additional tags.
func (c *Counter) With(labelValues ...string) harmonimetric.Counter {
	return &Counter{name: c.name, tags: append(c.tags, labelValues...)}
}

// Add increments the counter by delta.
func (c *Counter) Add(delta float64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value += delta
}

// Value returns the current counter value (for testing).
func (c *Counter) Value() float64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.value
}

// Gauge implements harmonimetric.Gauge for DogStatsD.
type Gauge struct {
	mu    sync.Mutex
	name  string
	tags  []string
	value float64
}

// NewGauge returns a Gauge that emits DogStatsD gauge metrics.
func NewGauge(name string, tags ...string) *Gauge {
	return &Gauge{name: name, tags: tags}
}

// With returns a new Gauge with additional tags.
func (g *Gauge) With(labelValues ...string) harmonimetric.Gauge {
	return &Gauge{name: g.name, tags: append(g.tags, labelValues...)}
}

// Set sets the gauge to an absolute value.
func (g *Gauge) Set(value float64) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.value = value
}

// Add increments or decrements the gauge by delta.
func (g *Gauge) Add(delta float64) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.value += delta
}

// Value returns the current gauge value (for testing).
func (g *Gauge) Value() float64 {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.value
}

// Histogram implements harmonimetric.Histogram for DogStatsD.
type Histogram struct {
	mu    sync.Mutex
	name  string
	tags  []string
	vals  []float64
}

// NewHistogram returns a Histogram that emits DogStatsD histogram metrics.
func NewHistogram(name string, tags ...string) *Histogram {
	return &Histogram{name: name, tags: tags}
}

// With returns a new Histogram with additional tags.
func (h *Histogram) With(labelValues ...string) harmonimetric.Histogram {
	return &Histogram{name: h.name, tags: append(h.tags, labelValues...)}
}

// Observe records a value in the histogram.
func (h *Histogram) Observe(value float64) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.vals = append(h.vals, value)
}

// Observations returns recorded values (for testing).
func (h *Histogram) Observations() []float64 {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.vals
}
