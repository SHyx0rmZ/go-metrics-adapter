package metrics_adapter

import (
	"github.com/golang/protobuf/proto"
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
)

type description prometheus.Desc

type gaugeAdapter struct {
	metric func() float64
	*description
}

type histogramAdapter struct {
	count      func() uint64
	sum        func() float64
	percentile func(p float64) uint64
	*description
}

func (d *description) Describe(ch chan<- *prometheus.Desc) {
	ch <- (*prometheus.Desc)(d)
}

func (d *description) Desc() *prometheus.Desc {
	return (*prometheus.Desc)(d)
}

func (a gaugeAdapter) Collect(ch chan<- prometheus.Metric) {
	ch <- a
}

func (a gaugeAdapter) Write(m *dto.Metric) error {
	m.Reset()
	m.Gauge = &dto.Gauge{
		Value: proto.Float64(a.metric()),
	}
	return nil
}

func (a histogramAdapter) Collect(ch chan<- prometheus.Metric) {
	ch <- a
}

func (a histogramAdapter) Write(m *dto.Metric) error {
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
