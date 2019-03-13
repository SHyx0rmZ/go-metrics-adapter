package metrics_adapter

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rcrowley/go-metrics"
)

type meter struct {
	metrics.Meter
	gaugeAdapter
}

func NewMeter(name string, metric metrics.Meter) interface {
	prometheus.Collector
	metrics.Meter
} {
	return meter{
		Meter: metric,
		gaugeAdapter: gaugeAdapter{
			metric: func(snapshot interface{}) float64 {
				return float64(snapshot.(metrics.Meter).Count())
			},
			snapshot: func() interface{} {
				return metric.Snapshot()
			},
			description: newDescriptionFrom(name),
		},
	}
}
