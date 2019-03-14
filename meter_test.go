package metrics_adapter

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"github.com/rcrowley/go-metrics"
	"testing"
)

func TestMeter(t *testing.T) {
	c := NewMeter("metric", metrics.NewMeter())
	m := make(chan prometheus.Metric, 1)

	for _, tt := range []struct {
		Mark int64
		Val  float64
	}{
		{
			Mark: 0,
			Val:  0,
		},
		{
			Mark: 2,
			Val:  2,
		},
		{
			Mark: 3,
			Val:  5,
		},
		{
			Mark: 0,
			Val:  5,
		},
	} {
		t.Run(fmt.Sprintf("Meter(%d,%d)", tt.Mark, int64(tt.Val)), func(t *testing.T) {
			c.Mark(tt.Mark)
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
