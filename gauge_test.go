package metrics_adapter

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"github.com/rcrowley/go-metrics"
	"testing"
)

func TestGauge(t *testing.T) {
	c := NewGauge("metric", metrics.NewGauge())
	m := make(chan prometheus.Metric, 1)

	for _, tt := range []struct {
		Update int64
		Val    float64
	}{
		{
			Update: 0,
			Val:    0,
		},
		{
			Update: 2,
			Val:    2,
		},
		{
			Update: 3,
			Val:    3,
		},
		{
			Update: 0,
			Val:    0,
		},
	} {
		t.Run(fmt.Sprintf("Gauge(%d,%d)", tt.Update, int64(tt.Val)), func(t *testing.T) {
			c.Update(tt.Update)
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

func TestGaugeFloat64(t *testing.T) {
	c := NewGaugeFloat64("metric", metrics.NewGaugeFloat64())
	m := make(chan prometheus.Metric, 1)

	for _, tt := range []struct {
		Update float64
		Val    float64
	}{
		{
			Update: 0,
			Val:    0,
		},
		{
			Update: 2,
			Val:    2,
		},
		{
			Update: 3,
			Val:    3,
		},
		{
			Update: 0,
			Val:    0,
		},
	} {
		t.Run(fmt.Sprintf("GaugeFloat64(%d,%d)", int64(tt.Update), int64(tt.Val)), func(t *testing.T) {
			c.Update(tt.Update)
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
