package metrics_adapter

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rcrowley/go-metrics"
)

type timer struct {
	metrics.Timer
	histogramAdapter
}

func NewTimer(name string, metric metrics.Timer) interface {
	prometheus.Collector
	metrics.Timer
} {
	return timer{
		Timer: metric,
		histogramAdapter: histogramAdapter{
			count: func() uint64 {
				return uint64(metric.Count())
			},
			sum: func() float64 {
				return float64(metric.Sum())
			},
			percentile: func(p float64) uint64 {
				return uint64(metric.Percentile(p))
			},
			description: newDescriptionFrom(name),
		},
	}
}
