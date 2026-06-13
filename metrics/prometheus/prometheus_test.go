package prometheus_test

import (
	"testing"

	harmonimetric "github.com/harmonikit/harmoni/metrics"
	prommetrics "github.com/harmonikit/kit/metrics/prometheus"
	"github.com/prometheus/client_golang/prometheus"
)

func TestCounter_Add(t *testing.T) {
	cv := prometheus.NewCounterVec(prometheus.CounterOpts{Name: "test_total"}, []string{"method"})
	c := prommetrics.NewCounter(cv)

	c.With("method", "GET").Add(1)
	// No panic — success.
}

func TestCounter_With(t *testing.T) {
	cv := prometheus.NewCounterVec(prometheus.CounterOpts{Name: "test2_total"}, []string{"method", "status"})
	c := prommetrics.NewCounter(cv)

	labeled := c.With("method", "POST").With("status", "200")
	labeled.Add(1)
}

func TestGauge_Set(t *testing.T) {
	gv := prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "test_gauge"}, []string{"region"})
	g := prommetrics.NewGauge(gv)

	g.With("region", "us-east").Set(100)
}

func TestGauge_Add(t *testing.T) {
	gv := prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "test_gauge2"}, []string{})
	g := prommetrics.NewGauge(gv)

	g.Add(5)
	g.Add(-2)
}

func TestHistogram_Observe(t *testing.T) {
	hv := prometheus.NewHistogramVec(prometheus.HistogramOpts{Name: "test_latency"}, []string{"path"})
	h := prommetrics.NewHistogram(hv)

	h.With("path", "/api/users").Observe(0.012)
}

func TestInterfaces(t *testing.T) {
	var _ harmonimetric.Counter = prommetrics.NewCounter(prometheus.NewCounterVec(prometheus.CounterOpts{Name: "x"}, nil))
	var _ harmonimetric.Gauge = prommetrics.NewGauge(prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "y"}, nil))
	var _ harmonimetric.Histogram = prommetrics.NewHistogram(prometheus.NewHistogramVec(prometheus.HistogramOpts{Name: "z"}, nil))
}