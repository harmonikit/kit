// Package prometheus adapts github.com/prometheus/client_golang to the
// harmoni/metrics interfaces.
package prometheus

import (
	harmonimetric "github.com/harmonikit/harmoni/metrics"
	prom "github.com/prometheus/client_golang/prometheus"
)

// Compile-time interface conformance.
var (
	_ harmonimetric.Counter   = (*Counter)(nil)
	_ harmonimetric.Gauge     = (*Gauge)(nil)
	_ harmonimetric.Histogram = (*Histogram)(nil)
)

// Counter implements harmonimetric.Counter backed by a Prometheus CounterVec.
type Counter struct {
	cv  *prom.CounterVec
	lvs []string
}

// NewCounter returns a Counter backed by a Prometheus CounterVec.
func NewCounter(cv *prom.CounterVec) *Counter {
	return &Counter{cv: cv, lvs: nil}
}

func (c *Counter) With(labelValues ...string) harmonimetric.Counter {
	return &Counter{cv: c.cv, lvs: append(c.lvs, labelValues...)}
}

func (c *Counter) Add(delta float64) {
	m, err := c.cv.GetMetricWithLabelValues(c.lvs...)
	if err != nil {
		return
	}
	m.Add(delta)
}

// Gauge implements harmonimetric.Gauge backed by a Prometheus GaugeVec.
type Gauge struct {
	gv  *prom.GaugeVec
	lvs []string
}

// NewGauge returns a Gauge backed by a Prometheus GaugeVec.
func NewGauge(gv *prom.GaugeVec) *Gauge {
	return &Gauge{gv: gv, lvs: nil}
}

func (g *Gauge) With(labelValues ...string) harmonimetric.Gauge {
	return &Gauge{gv: g.gv, lvs: append(g.lvs, labelValues...)}
}

func (g *Gauge) Set(value float64) {
	m, err := g.gv.GetMetricWithLabelValues(g.lvs...)
	if err != nil {
		return
	}
	m.Set(value)
}

func (g *Gauge) Add(delta float64) {
	m, err := g.gv.GetMetricWithLabelValues(g.lvs...)
	if err != nil {
		return
	}
	m.Add(delta)
}

// Histogram implements harmonimetric.Histogram backed by a Prometheus HistogramVec.
type Histogram struct {
	hv  *prom.HistogramVec
	lvs []string
}

// NewHistogram returns a Histogram backed by a Prometheus HistogramVec.
func NewHistogram(hv *prom.HistogramVec) *Histogram {
	return &Histogram{hv: hv, lvs: nil}
}

func (h *Histogram) With(labelValues ...string) harmonimetric.Histogram {
	return &Histogram{hv: h.hv, lvs: append(h.lvs, labelValues...)}
}

func (h *Histogram) Observe(value float64) {
	m, err := h.hv.GetMetricWithLabelValues(h.lvs...)
	if err != nil {
		return
	}
	m.Observe(value)
}
