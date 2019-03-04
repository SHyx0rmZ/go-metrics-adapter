package metrics_adapter

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rcrowley/go-metrics"
)

type healthcheckAdapter struct {
	metrics.Healthcheck
	__gaugeAdapter
}

func NewHealthcheckAdapter(s string, m metrics.Healthcheck) interface {
	prometheus.Collector
	metrics.Healthcheck
} {
	return healthcheckAdapter{
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
