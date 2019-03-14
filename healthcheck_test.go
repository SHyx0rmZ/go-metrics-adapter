package metrics_adapter

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"github.com/rcrowley/go-metrics"
	"testing"
)

func TestHealthcheck(t *testing.T) {
	e := make(chan error, 1)
	c := NewHealthcheck("metric", metrics.NewHealthcheck(func(h metrics.Healthcheck) {
		err := <-e
		if err != nil {
			h.Unhealthy(err)
		} else {
			h.Healthy()
		}
	}))
	m := make(chan prometheus.Metric, 1)

	for _, tt := range []struct {
		Err error
		Val float64
	}{
		{
			Err: nil,
			Val: 1,
		},
		{
			Err: fmt.Errorf("error"),
			Val: 0,
		},
	} {
		t.Run(fmt.Sprintf("HealthCheck(%v,%d)", tt.Err, int64(tt.Val)), func(t *testing.T) {
			e <- tt.Err
			c.Check()
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
