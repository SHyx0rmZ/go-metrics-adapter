package metrics_adapter

import (
	"github.com/golang/protobuf/proto"
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
)

type __gaugeAdapter struct {
	metric func() float64
	desc   *prometheus.Desc
}

func (a __gaugeAdapter) Describe(ch chan<- *prometheus.Desc) {
	ch <- a.desc
}

func (a __gaugeAdapter) Desc() *prometheus.Desc {
	return a.desc
}

func (a __gaugeAdapter) Collect(ch chan<- prometheus.Metric) {
	ch <- a
}

func (a __gaugeAdapter) Write(m *dto.Metric) error {
	m.Reset()
	m.Gauge = &dto.Gauge{
		Value: proto.Float64(a.metric()),
	}
	return nil
}

type __histogramAdapter struct {
	count      func() uint64
	sum        func() float64
	percentile func(p float64) uint64
	desc       *prometheus.Desc
}

func (a __histogramAdapter) Describe(ch chan<- *prometheus.Desc) {
	ch <- a.desc
}

func (a __histogramAdapter) Desc() *prometheus.Desc {
	return a.desc
}

func (a __histogramAdapter) Collect(ch chan<- prometheus.Metric) {
	ch <- a
}

func (a __histogramAdapter) Write(m *dto.Metric) error {
	m.Reset()
	m.Histogram = &dto.Histogram{
		SampleCount: proto.Uint64(a.count()),
		SampleSum:   proto.Float64(a.sum()),
	}
	for _, b := range prometheus.DefBuckets {
		m.Histogram.Bucket = append(m.Histogram.Bucket, &dto.Bucket{
			CumulativeCount: proto.Uint64(uint64(a.percentile(b))),
			UpperBound:      proto.Float64(b),
		})
	}
	return nil
}
