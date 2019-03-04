package metrics_adapter

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rcrowley/go-metrics"
)

type meterAdapter struct {
	metrics.Meter
	__gaugeAdapter
}

func NewMeterAdapter(s string, m metrics.Meter) interface {
	prometheus.Collector
	metrics.Meter
} {
	return meterAdapter{
		Meter: m,
		__gaugeAdapter: __gaugeAdapter{
			metric: intToFloat(m.Count),
			desc:   desc(s),
		},
	}
}
