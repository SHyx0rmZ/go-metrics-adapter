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
			count: func(snapshot interface{}) uint64 {
				return uint64(snapshot.(metrics.Timer).Count())
			},
			sum: func(snapshot interface{}) float64 {
				return float64(snapshot.(metrics.Timer).Sum())
			},
			percentile: func(snapshot interface{}, p float64) uint64 {
				return uint64(snapshot.(metrics.Timer).Percentile(p))
			},
			snapshot: func() interface{} {
				return metric.Snapshot()
			},
			description: newDescriptionFrom(name),
		},
	}
}
