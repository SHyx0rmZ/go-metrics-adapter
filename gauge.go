package metrics_adapter

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rcrowley/go-metrics"
)

type gauge struct {
	metrics.Gauge
	gaugeAdapter
}

func NewGauge(name string, metric metrics.Gauge) interface {
	prometheus.Collector
	metrics.Gauge
} {
	return gauge{
		Gauge: metric,
		gaugeAdapter: gaugeAdapter{
			metric: func() float64 {
				return float64(metric.Value())
			},
			description: newDescriptionFrom(name),
		},
	}
}

type gaugeFloat64 struct {
	metrics.GaugeFloat64
	gaugeAdapter
}

func NewGaugeFloat64(name string, metric metrics.GaugeFloat64) interface {
	prometheus.Collector
	metrics.GaugeFloat64
} {
	return gaugeFloat64{
		GaugeFloat64: metric,
		gaugeAdapter: gaugeAdapter{
			metric:      metric.Value,
			description: newDescriptionFrom(name),
		},
	}
}
