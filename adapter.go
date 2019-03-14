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

type gaugeAdapterMetric struct {
	metric   func(snapshot interface{}) float64
	snapshot interface{}
	*description
}

func (a gaugeAdapter) Collect(ch chan<- prometheus.Metric) {
	ch <- gaugeAdapterMetric{
		snapshot:    a.snapshot(),
		metric:      a.metric,
		description: a.description,
	}
}

func (a gaugeAdapterMetric) Write(m *dto.Metric) error {
	m.Reset()
	m.Gauge = &dto.Gauge{
		Value: proto.Float64(a.metric(a.snapshot)),
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

type summaryAdapterMetric struct {
	count      func(snapshot interface{}) uint64
	sum        func(snapshot interface{}) float64
	percentile func(snapshot interface{}, p float64) float64
	snapshot   interface{}
	*description
}

func (a summaryAdapter) Collect(ch chan<- prometheus.Metric) {
	ch <- summaryAdapterMetric{
		count:       a.count,
		sum:         a.sum,
		percentile:  a.percentile,
		snapshot:    a.snapshot(),
		description: a.description,
	}
}

func (a summaryAdapterMetric) Write(m *dto.Metric) error {
	m.Reset()
	m.Summary = &dto.Summary{
		SampleCount: proto.Uint64(a.count(a.snapshot)),
		SampleSum:   proto.Float64(a.sum(a.snapshot)),
	}
	for _, p := range quantiles {
		m.Summary.Quantile = append(m.Summary.Quantile, &dto.Quantile{
			Quantile: proto.Float64(p),
			Value:    proto.Float64(a.percentile(a.snapshot, p)),
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
