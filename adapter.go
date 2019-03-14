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

// Collect is called by the Prometheus registry when collecting
// metrics. The implementation sends each collected metric via the
// provided channel and returns once the last metric has been sent. The
// descriptor of each sent metric is one of those returned by Describe
// (unless the Collector is unchecked, see above). Returned metrics that
// share the same descriptor must differ in their variable label
// values.
//
// This method may be called concurrently and must therefore be
// implemented in a concurrency safe way. Blocking occurs at the expense
// of total performance of rendering all registered metrics. Ideally,
// Collector implementations support concurrent readers.
func (a gaugeAdapter) Collect(ch chan<- prometheus.Metric) {
	ch <- gaugeAdapterMetric{
		snapshot:    a.snapshot(),
		metric:      a.metric,
		description: a.description,
	}
}

// Write encodes the Metric into a "Metric" Protocol Buffer data
// transmission object.
//
// Metric implementations must observe concurrency safety as reads of
// this metric may occur at any time, and any blocking occurs at the
// expense of total performance of rendering all registered
// metrics. Ideally, Metric implementations should support concurrent
// readers.
//
// While populating dto.Metric, it is the responsibility of the
// implementation to ensure validity of the Metric protobuf (like valid
// UTF-8 strings or syntactically valid metric and label names). It is
// recommended to sort labels lexicographically. Callers of Write should
// still make sure of sorting if they depend on it.
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

// Collect is called by the Prometheus registry when collecting
// metrics. The implementation sends each collected metric via the
// provided channel and returns once the last metric has been sent. The
// descriptor of each sent metric is one of those returned by Describe
// (unless the Collector is unchecked, see above). Returned metrics that
// share the same descriptor must differ in their variable label
// values.
//
// This method may be called concurrently and must therefore be
// implemented in a concurrency safe way. Blocking occurs at the expense
// of total performance of rendering all registered metrics. Ideally,
// Collector implementations support concurrent readers.
func (a summaryAdapter) Collect(ch chan<- prometheus.Metric) {
	ch <- summaryAdapterMetric{
		count:       a.count,
		sum:         a.sum,
		percentile:  a.percentile,
		snapshot:    a.snapshot(),
		description: a.description,
	}
}

// Write encodes the Metric into a "Metric" Protocol Buffer data
// transmission object.
//
// Metric implementations must observe concurrency safety as reads of
// this metric may occur at any time, and any blocking occurs at the
// expense of total performance of rendering all registered
// metrics. Ideally, Metric implementations should support concurrent
// readers.
//
// While populating dto.Metric, it is the responsibility of the
// implementation to ensure validity of the Metric protobuf (like valid
// UTF-8 strings or syntactically valid metric and label names). It is
// recommended to sort labels lexicographically. Callers of Write should
// still make sure of sorting if they depend on it.
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
