package metrics_adapter

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rcrowley/go-metrics"
)

type gauge struct {
	metrics.Gauge
	gaugeAdapter
}

func NewGauge(name string, m metrics.Gauge) interface {
	prometheus.Collector
	metrics.Gauge
} {
	return gauge{
		Gauge: m,
		gaugeAdapter: gaugeAdapter{
			metric:      intToFloat(m.Value),
			description: newDescriptionFrom(name),
		},
	}
}

type gaugeFloat64 struct {
	metrics.GaugeFloat64
	gaugeAdapter
}

func NewGaugeFloat64(name string, m metrics.GaugeFloat64) interface {
	prometheus.Collector
	metrics.GaugeFloat64
} {
	return gaugeFloat64{
		GaugeFloat64: m,
		gaugeAdapter: gaugeAdapter{
			metric:      m.Value,
			description: newDescriptionFrom(name),
		},
	}
}
