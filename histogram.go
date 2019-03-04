package metrics_adapter

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rcrowley/go-metrics"
)

type histogramAdapter struct {
	metrics.Histogram
	__histogramAdapter
}

func NewHistogramAdapter(s string, m metrics.Histogram) interface {
	prometheus.Collector
	metrics.Histogram
} {
	return histogramAdapter{
		Histogram: m,
		__histogramAdapter: __histogramAdapter{
			count:      intToUint(m.Count),
			sum:        intToFloat(m.Sum),
			percentile: floatToUint(m.Percentile),
			desc:       desc(s),
		},
	}
}
