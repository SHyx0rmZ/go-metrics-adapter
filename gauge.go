package metrics_adapter

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rcrowley/go-metrics"
)

type gauge struct {
	metrics.Gauge
	__gaugeAdapter
}

func NewGauge(s string, m metrics.Gauge) interface {
	prometheus.Collector
	metrics.Gauge
} {
	return gauge{
		Gauge: m,
		__gaugeAdapter: __gaugeAdapter{
			metric: intToFloat(m.Value),
			desc:   desc(s),
		},
	}
}

type gaugeFloat64 struct {
	metrics.GaugeFloat64
	__gaugeAdapter
}

func NewGaugeFloat64(s string, m metrics.GaugeFloat64) interface {
	prometheus.Collector
	metrics.GaugeFloat64
} {
	return gaugeFloat64{
		GaugeFloat64: m,
		__gaugeAdapter: __gaugeAdapter{
			metric: m.Value,
			desc:   desc(s),
		},
	}
}
