package metrics_adapter

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rcrowley/go-metrics"
)

type gaugeAdapter struct {
	metrics.Gauge
	__gaugeAdapter
}

func NewGaugeAdapter(s string, m metrics.Gauge) interface {
	prometheus.Collector
	metrics.Gauge
} {
	return gaugeAdapter{
		Gauge: m,
		__gaugeAdapter: __gaugeAdapter{
			metric: intToFloat(m.Value),
			desc:   desc(s),
		},
	}
}

type gaugeFloat64Adapter struct {
	metrics.GaugeFloat64
	__gaugeAdapter
}

func NewGaugeFloat64Adapter(s string, m metrics.GaugeFloat64) interface {
	prometheus.Collector
	metrics.GaugeFloat64
} {
	return gaugeFloat64Adapter{
		GaugeFloat64: m,
		__gaugeAdapter: __gaugeAdapter{
			metric: m.Value,
			desc:   desc(s),
		},
	}
}
