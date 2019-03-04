package metrics_adapter

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rcrowley/go-metrics"
)

type meter struct {
	metrics.Meter
	gaugeAdapter
}

func NewMeter(s string, m metrics.Meter) interface {
	prometheus.Collector
	metrics.Meter
} {
	return meter{
		Meter: m,
		gaugeAdapter: gaugeAdapter{
			metric: intToFloat(m.Count),
			desc:   desc(s),
		},
	}
}
