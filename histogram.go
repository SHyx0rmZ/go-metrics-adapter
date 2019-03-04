package metrics_adapter

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rcrowley/go-metrics"
)

type histogram struct {
	metrics.Histogram
	histogramAdapter
}

func NewHistogram(name string, metric metrics.Histogram) interface {
	prometheus.Collector
	metrics.Histogram
} {
	return histogram{
		Histogram: metric,
		histogramAdapter: histogramAdapter{
			count:       intToUint(metric.Count),
			sum:         intToFloat(metric.Sum),
			percentile:  floatToUint(metric.Percentile),
			description: newDescriptionFrom(name),
		},
	}
}
