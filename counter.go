package metrics_adapter

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rcrowley/go-metrics"
)

type counter struct {
	metrics.Counter
	gaugeAdapter
}

// NewCounter turns metric into a prometheus.Collector. The description will
// be taken from name.
func NewCounter(name string, metric metrics.Counter) interface {
	prometheus.Collector
	metrics.Counter
} {
	return counter{
		Counter: metric,
		gaugeAdapter: gaugeAdapter{
			metric: func(snapshot interface{}) float64 {
				return float64(snapshot.(metrics.Counter).Count())
			},
			snapshot: func() interface{} {
				return metric.Snapshot()
			},
			description: newDescriptionFrom(name),
		},
	}
}
