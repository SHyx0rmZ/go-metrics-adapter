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
