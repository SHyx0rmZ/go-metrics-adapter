package metrics_adapter

import (
	"github.com/prometheus/client_golang/prometheus"
	"strings"
)

type description prometheus.Desc

func (d *description) Describe(ch chan<- *prometheus.Desc) {
	ch <- (*prometheus.Desc)(d)
}

func (d *description) Desc() *prometheus.Desc {
	return (*prometheus.Desc)(d)
}

func newDescriptionFrom(name string) *description {
	name = strings.Replace(name, "-", "_", -1)

	return (*description)(prometheus.NewDesc(
		"witches_io_metrics_adapter_"+name,
		"Adapter for "+name,
		nil,
		nil,
	))
}
