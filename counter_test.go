package metrics_adapter

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"github.com/rcrowley/go-metrics"
	"testing"
)

func TestCounter(t *testing.T) {
	c := NewCounter("metric", metrics.NewCounter())
	m := make(chan prometheus.Metric, 1)

	for _, tt := range []struct {
		Inc int64
		Val float64
	}{
		{
			Inc: 0,
			Val: 0,
		},
		{
			Inc: 2,
			Val: 2,
		},
		{
			Inc: 3,
			Val: 5,
		},
		{
			Inc: 0,
			Val: 5,
		},
	} {
		t.Run(fmt.Sprintf("Counter(%d,%d)", tt.Inc, int64(tt.Val)), func(t *testing.T) {
			c.Inc(tt.Inc)
			c.Collect(m)
			p := <-m
			var d dto.Metric
			err := p.Write(&d)
			if err != nil {
				t.Error(err)
			}
			if d.Gauge == nil {
				t.Fatal("Gauge == nil")
			}
			if *d.Gauge.Value != tt.Val {
				t.Errorf("Gauge.Value: got %v, want %v", *d.Gauge.Value, tt.Val)
			}
		})
	}
}
