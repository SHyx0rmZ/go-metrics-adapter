package metrics_adapter

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rcrowley/go-metrics"
)

type counterAdapter struct {
	metrics.Counter
	__gaugeAdapter
}

func NewCounterAdapter(s string, m metrics.Counter) interface {
	prometheus.Collector
	metrics.Counter
} {
	return counterAdapter{
		Counter: m,
		__gaugeAdapter: __gaugeAdapter{
			metric: intToFloat(m.Count),
			desc:   desc(s),
		},
	}
}
