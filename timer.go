package metrics_adapter

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rcrowley/go-metrics"
)

type timer struct {
	metrics.Timer
	__histogramAdapter
}

func NewTimer(s string, m metrics.Timer) interface {
	prometheus.Collector
	metrics.Timer
} {
	return timer{
		Timer: m,
		__histogramAdapter: __histogramAdapter{
			count:      intToUint(m.Count),
			sum:        intToFloat(m.Sum),
			percentile: floatToUint(m.Percentile),
			desc:       desc(s),
		},
	}
}
