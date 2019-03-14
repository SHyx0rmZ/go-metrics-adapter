package metrics_adapter

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rcrowley/go-metrics"
)

type timer struct {
	metrics.Timer
	summaryAdapter
}

// NewTimer turns metric into a prometheus.Collector. The description will
// be taken from name.
func NewTimer(name string, metric metrics.Timer) interface {
	prometheus.Collector
	metrics.Timer
} {
	return timer{
		Timer: metric,
		summaryAdapter: summaryAdapter{
			count: func(snapshot interface{}) uint64 {
				return uint64(snapshot.(metrics.Timer).Count())
			},
			sum: func(snapshot interface{}) float64 {
				return float64(snapshot.(metrics.Timer).Sum())
			},
			percentile: func(snapshot interface{}, p float64) float64 {
				return snapshot.(metrics.Timer).Percentile(p)
			},
			snapshot: func() interface{} {
				return metric.Snapshot()
			},
			description: newDescriptionFrom(name),
		},
	}
}
