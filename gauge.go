package metrics_adapter

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rcrowley/go-metrics"
)

type gauge struct {
	metrics.Gauge
	gaugeAdapter
}

// NewGauge turns metric into a prometheus.Collector. The description will
// be taken from name.
func NewGauge(name string, metric metrics.Gauge) interface {
	prometheus.Collector
	metrics.Gauge
} {
	return gauge{
		Gauge: metric,
		gaugeAdapter: gaugeAdapter{
			metric: func(snapshot interface{}) float64 {
				return float64(snapshot.(metrics.Gauge).Value())
			},
			snapshot: func() interface{} {
				return metric.Snapshot()
			},
			description: newDescriptionFrom(name),
		},
	}
}

type gaugeFloat64 struct {
	metrics.GaugeFloat64
	gaugeAdapter
}

// NewGaugeFloat64 turns metric into a prometheus.Collector. The description
// will be taken from name.
func NewGaugeFloat64(name string, metric metrics.GaugeFloat64) interface {
	prometheus.Collector
	metrics.GaugeFloat64
} {
	return gaugeFloat64{
		GaugeFloat64: metric,
		gaugeAdapter: gaugeAdapter{
			metric: func(snapshot interface{}) float64 {
				return snapshot.(metrics.GaugeFloat64).Value()
			},
			snapshot: func() interface{} {
				return metric.Snapshot()
			},
			description: newDescriptionFrom(name),
		},
	}
}
