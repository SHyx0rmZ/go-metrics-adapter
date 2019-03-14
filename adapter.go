package metrics_adapter

import (
	"github.com/golang/protobuf/proto"
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"sort"
)

type gaugeAdapter struct {
	metric   func(snapshot interface{}) float64
	snapshot func() interface{}
	*description
}

func (a gaugeAdapter) Collect(ch chan<- prometheus.Metric) {
	ch <- a
}

func (a gaugeAdapter) Write(m *dto.Metric) error {
	m.Reset()
	s := a.snapshot()
	m.Gauge = &dto.Gauge{
		Value: proto.Float64(a.metric(s)),
	}
	return nil
}

type summaryAdapter struct {
	count      func(snapshot interface{}) uint64
	sum        func(snapshot interface{}) float64
	percentile func(snapshot interface{}, p float64) float64
	snapshot   func() interface{}
	*description
}

func (a summaryAdapter) Collect(ch chan<- prometheus.Metric) {
	ch <- a
}

func (a summaryAdapter) Write(m *dto.Metric) error {
	m.Reset()
	s := a.snapshot()
	m.Summary = &dto.Summary{
		SampleCount: proto.Uint64(a.count(s)),
		SampleSum:   proto.Float64(a.sum(s)),
	}
	for _, p := range quantiles {
		m.Summary.Quantile = append(m.Summary.Quantile, &dto.Quantile{
			Quantile: proto.Float64(p),
			Value:    proto.Float64(a.percentile(s, p)),
		})
	}
	return nil
}

var quantiles []float64

func init() {
	for q := range prometheus.DefObjectives {
		quantiles = append(quantiles, q)
	}
	sort.Float64s(quantiles)
}
