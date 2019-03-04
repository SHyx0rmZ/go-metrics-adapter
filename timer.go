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
			count:       intToUint(metric.Count),
			sum:         intToFloat(metric.Sum),
			percentile:  floatToUint(metric.Percentile),
			description: newDescriptionFrom(name),
		},
	}
}
