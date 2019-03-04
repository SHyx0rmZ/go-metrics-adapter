package metrics_adapter

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rcrowley/go-metrics"
)

type counter struct {
	metrics.Counter
	gaugeAdapter
}

func NewCounter(name string, metric metrics.Counter) interface {
	prometheus.Collector
	metrics.Counter
} {
	return counter{
		Counter: metric,
		gaugeAdapter: gaugeAdapter{
			metric:      intToFloat(metric.Count),
			description: newDescriptionFrom(name),
		},
	}
}
