package metrics_adapter

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"github.com/rcrowley/go-metrics"
	"testing"
)

func TestHistogram(t *testing.T) {
	type bucket struct {
		Quantile float64
		Value    float64
	}

	c := NewHistogram("metric", metrics.NewHistogram(metrics.NewUniformSample(16)))
	m := make(chan prometheus.Metric, 1)

	for _, tt := range []struct {
		Update    int64
		Sum       int64
		Count     uint64
		Quantiles []bucket
	}{
		{
			Update: 0,
			Sum:    0,
			Count:  1,
			Quantiles: []bucket{
				{
					Quantile: 0.5,
					Value:    0,
				},
				{
					Quantile: 0.9,
					Value:    0,
				},
				{
					Quantile: 0.99,
					Value:    0,
				},
			},
		},
		{
			Update: 2,
			Sum:    2,
			Count:  2,
			Quantiles: []bucket{
				{
					Quantile: 0.5,
					Value:    1,
				},
				{
					Quantile: 0.9,
					Value:    2,
				},
				{
					Quantile: 0.99,
					Value:    2,
				},
			},
		},
		{
			Update: 3,
			Sum:    5,
			Count:  3,
			Quantiles: []bucket{
				{
					Quantile: 0.5,
					Value:    2,
				},
				{
					Quantile: 0.9,
					Value:    3,
				},
				{
					Quantile: 0.99,
					Value:    3,
				},
			},
		},
	} {
		t.Run(fmt.Sprintf("Histogram(%d,%d,%d)", int64(tt.Update), tt.Sum, int64(tt.Count)), func(t *testing.T) {
			c.Update(tt.Update)
			c.Collect(m)
			p := <-m
			var d dto.Metric
			err := p.Write(&d)
			if err != nil {
				t.Error(err)
			}
			if d.Gauge != nil {
				t.Fatal("Gauge != nil")
			}
			if d.Summary == nil {
				t.Fatal("Summary == nil")
			}
			if *d.Summary.SampleCount != tt.Count {
				t.Errorf("Summary.SampleCount: got %v, want %v", *d.Summary.SampleCount, tt.Count)
			}
			if *d.Summary.SampleSum != float64(tt.Sum) {
				t.Errorf("Summary.SampleSum: got %v, want %v", *d.Summary.SampleSum, tt.Sum)
			}
			if d.Summary.Quantile == nil {
				t.Fatal("Summary.Quantile == nil")
			}
			for i, b := range tt.Quantiles {
				if *d.Summary.Quantile[i].Quantile != b.Quantile {
					t.Errorf("Summary.Quantile[%d].Quantile: got %v, want %v", i, *d.Summary.Quantile[i].Quantile, b.Quantile)
				}
				if *d.Summary.Quantile[i].Value != b.Value {
					t.Errorf("Summary.Quantile[%d].Value: got %v, want %v", i, *d.Summary.Quantile[i].Value, b.Value)
				}
			}
		})
	}
}
