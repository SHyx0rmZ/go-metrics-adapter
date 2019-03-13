package metrics_adapter

import (
	"github.com/golang/protobuf/proto"
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
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

type histogramAdapter struct {
	count      func(snapshot interface{}) uint64
	sum        func(snapshot interface{}) float64
	percentile func(snapshot interface{}, p float64) uint64
	snapshot   func() interface{}
	*description
}

func (a histogramAdapter) Collect(ch chan<- prometheus.Metric) {
	ch <- a
}

func (a histogramAdapter) Write(m *dto.Metric) error {
	m.Reset()
	s := a.snapshot()
	m.Histogram = &dto.Histogram{
		SampleCount: proto.Uint64(a.count(s)),
		SampleSum:   proto.Float64(a.sum(s)),
	}
	for _, b := range prometheus.DefBuckets {
		m.Histogram.Bucket = append(m.Histogram.Bucket, &dto.Bucket{
			CumulativeCount: proto.Uint64(a.percentile(s, b)),
			UpperBound:      proto.Float64(b),
		})
	}
	return nil
}
