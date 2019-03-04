package metrics_adapter

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rcrowley/go-metrics"
)

type counter struct {
	metrics.Counter
	__gaugeAdapter
}

func NewCounter(s string, m metrics.Counter) interface {
	prometheus.Collector
	metrics.Counter
} {
	return counter{
		Counter: m,
		__gaugeAdapter: __gaugeAdapter{
			metric: intToFloat(m.Count),
			desc:   desc(s),
		},
	}
}
