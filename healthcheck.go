package metrics_adapter

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rcrowley/go-metrics"
)

type healthcheck struct {
	metrics.Healthcheck
	gaugeAdapter
}

func NewHealthcheck(name string, metric metrics.Healthcheck) interface {
	prometheus.Collector
	metrics.Healthcheck
} {
	return healthcheck{
		Healthcheck: metric,
		gaugeAdapter: gaugeAdapter{
			metric: func() float64 {
				if metric.Error() != nil {
					return 0
				}
				return 1
			},
			description: newDescriptionFrom(name),
		},
	}
}
