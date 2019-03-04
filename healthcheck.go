package metrics_adapter

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rcrowley/go-metrics"
)

type healthcheck struct {
	metrics.Healthcheck
	__gaugeAdapter
}

func NewHealthcheck(s string, m metrics.Healthcheck) interface {
	prometheus.Collector
	metrics.Healthcheck
} {
	return healthcheck{
		Healthcheck: m,
		__gaugeAdapter: __gaugeAdapter{
			metric: func() float64 {
				if m.Error() != nil {
					return 0
				}
				return 1
			},
			desc: desc(s),
		},
	}
}
