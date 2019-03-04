package metrics_adapter

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rcrowley/go-metrics"
)

type timer struct {
	metrics.Timer
	histogramAdapter
}

func NewTimer(s string, m metrics.Timer) interface {
	prometheus.Collector
	metrics.Timer
} {
	return timer{
		Timer: m,
		histogramAdapter: histogramAdapter{
			count:       intToUint(m.Count),
			sum:         intToFloat(m.Sum),
			percentile:  floatToUint(m.Percentile),
			description: newDescriptionFrom(s),
		},
	}
}
