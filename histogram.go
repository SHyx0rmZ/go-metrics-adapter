package metrics_adapter

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rcrowley/go-metrics"
)

type histogram struct {
	metrics.Histogram
	summaryAdapter
}

// NewHistogram turns metric into a prometheus.Collector. The description will
// be taken from name.
func NewHistogram(name string, metric metrics.Histogram) interface {
	prometheus.Collector
	metrics.Histogram
} {
	return histogram{
		Histogram: metric,
		summaryAdapter: summaryAdapter{
			count: func(snapshot interface{}) uint64 {
				return uint64(snapshot.(metrics.Histogram).Count())
			},
			sum: func(snapshot interface{}) float64 {
				return float64(snapshot.(metrics.Histogram).Sum())
			},
			percentile: func(snapshot interface{}, p float64) float64 {
				return snapshot.(metrics.Histogram).Percentile(p)
			},
			snapshot: func() interface{} {
				return metric.Snapshot()
			},
			description: newDescriptionFrom(name),
		},
	}
}
