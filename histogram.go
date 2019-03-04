package metrics_adapter

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rcrowley/go-metrics"
)

type histogram struct {
	metrics.Histogram
	histogramAdapter
}

func NewHistogram(s string, m metrics.Histogram) interface {
	prometheus.Collector
	metrics.Histogram
} {
	return histogram{
		Histogram: m,
		histogramAdapter: histogramAdapter{
			count:       intToUint(m.Count),
			sum:         intToFloat(m.Sum),
			percentile:  floatToUint(m.Percentile),
			description: desc(s),
		},
	}
}
