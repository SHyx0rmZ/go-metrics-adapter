package metrics_adapter

import (
	"github.com/prometheus/client_golang/prometheus"
	"strings"
)

func desc(s string) *prometheus.Desc {
	return prometheus.NewDesc("witches_io_metrics_adapter_"+strings.Replace(s, "-", "_", -1), "Adapter for "+s, nil, nil)
}

func intToFloat(f func() int64) func() float64 {
	return func() float64 {
		return float64(f())
	}
}

func intToUint(f func() int64) func() uint64 {
	return func() uint64 {
		return uint64(f())
	}
}

func floatToUint(f func(float64) float64) func(float64) uint64 {
	return func(v float64) uint64 {
		return uint64(f(v))
	}
}
