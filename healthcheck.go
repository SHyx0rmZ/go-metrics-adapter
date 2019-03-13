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
			metric: func(snapshot interface{}) float64 {
				if snapshot.(metrics.Healthcheck).Error() != nil {
					return 0
				}
				return 1
			},
			snapshot: func() interface{} {
				return metric
			},
			description: newDescriptionFrom(name),
		},
	}
}
