package metrics_adapter

import (
	"reflect"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/rcrowley/go-metrics"
)

type registry struct {
	registerer prometheus.Registerer

	names map[string]prometheus.Collector
	mu    sync.Mutex
}

// NewRegistry turns registerer into a metrics.Registry.
func NewRegistry(registerer prometheus.Registerer) metrics.Registry {
	return &registry{
		registerer: registerer,
		names:      make(map[string]prometheus.Collector),
	}
}

// Call the given function for each registered metric.
func (a *registry) Each(f func(string, interface{})) {
	a.mu.Lock()
	defer a.mu.Unlock()
	for name, collector := range a.names {
		f(name, collector)
	}
}

// Get the metric by the given name or nil if none is registered.
func (a *registry) Get(name string) interface{} {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.names[name]
}

// GetAll metrics in the Registry.
func (a *registry) GetAll() map[string]map[string]interface{} {
	data := make(map[string]map[string]interface{})
	a.Each(func(name string, i interface{}) {
		values := make(map[string]interface{})
		switch metric := i.(type) {
		case metrics.Counter:
			values["count"] = metric.Count()
		case metrics.Gauge:
			values["value"] = metric.Value()
		case metrics.GaugeFloat64:
			values["value"] = metric.Value()
		case metrics.Healthcheck:
			values["error"] = nil
			metric.Check()
			if err := metric.Error(); nil != err {
				values["error"] = metric.Error().Error()
			}
		case metrics.Histogram:
			h := metric.Snapshot()
			ps := h.Percentiles([]float64{0.5, 0.75, 0.95, 0.99, 0.999})
			values["count"] = h.Count()
			values["min"] = h.Min()
			values["max"] = h.Max()
			values["mean"] = h.Mean()
			values["stddev"] = h.StdDev()
			values["median"] = ps[0]
			values["75%"] = ps[1]
			values["95%"] = ps[2]
			values["99%"] = ps[3]
			values["99.9%"] = ps[4]
		case metrics.Meter:
			m := metric.Snapshot()
			values["count"] = m.Count()
			values["1m.rate"] = m.Rate1()
			values["5m.rate"] = m.Rate5()
			values["15m.rate"] = m.Rate15()
			values["mean.rate"] = m.RateMean()
		case metrics.Timer:
			t := metric.Snapshot()
			ps := t.Percentiles([]float64{0.5, 0.75, 0.95, 0.99, 0.999})
			values["count"] = t.Count()
			values["min"] = t.Min()
			values["max"] = t.Max()
			values["mean"] = t.Mean()
			values["stddev"] = t.StdDev()
			values["median"] = ps[0]
			values["75%"] = ps[1]
			values["95%"] = ps[2]
			values["99%"] = ps[3]
			values["99.9%"] = ps[4]
			values["1m.rate"] = t.Rate1()
			values["5m.rate"] = t.Rate5()
			values["15m.rate"] = t.Rate15()
			values["mean.rate"] = t.RateMean()
		}
		data[name] = values
	})
	return data
}

// Gets an existing metric or registers the given one.
// The interface can be the metric to register if not found in registry,
// or a function returning the metric for lazy instantiation.
func (a *registry) GetOrRegister(name string, v interface{}) interface{} {
	c := a.Get(name)
	if c == nil {
		err := a.Register(name, v)
		if err != nil {
			panic(err)
		}
		return a.Get(name)
	}
	return c
}

// Register the given metric under the given name.
func (a *registry) Register(name string, v interface{}) error {
	c, ok := v.(prometheus.Collector)
	if !ok {
		if reflect.TypeOf(v).Kind() == reflect.Func {
			v = reflect.ValueOf(v).Call(nil)[0].Interface()
		}
		switch metric := v.(type) {
		case metrics.Counter:
			c = NewCounter(name, metric)
		case metrics.Gauge:
			c = NewGauge(name, metric)
		case metrics.GaugeFloat64:
			c = NewGaugeFloat64(name, metric)
		case metrics.Healthcheck:
			c = NewHealthcheck(name, metric)
		case metrics.Histogram:
			c = NewHistogram(name, metric)
		case metrics.Meter:
			c = NewMeter(name, metric)
		case metrics.Timer:
			c = NewTimer(name, metric)
		default:
			return ErrExpectedCollector
		}
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	err := a.registerer.Register(c)
	if err != nil {
		return err
	}
	a.names[name] = c
	return nil
}

// Run all registered healthchecks.
func (a *registry) RunHealthchecks() {
	a.Each(func(name string, metric interface{}) {
		if h, ok := metric.(metrics.Healthcheck); ok {
			h.Check()
		}
	})
}

// Unregister the metric with the given name.
func (a *registry) Unregister(name string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	metric, ok := a.names[name]
	if !ok {
		return
	}
	if i, ok := metric.(metrics.Stoppable); ok {
		i.Stop()
	}
	a.registerer.Unregister(metric)
	delete(a.names, name)
}

// Unregister all metrics.  (Mostly for testing.)
func (a *registry) UnregisterAll() {
	a.mu.Lock()
	defer a.mu.Unlock()
	for name, c := range a.names {
		metric, ok := a.names[name]
		if !ok {
			continue
		}
		if i, ok := metric.(metrics.Stoppable); ok {
			i.Stop()
		}
		a.registerer.Unregister(c)
		delete(a.names, name)
	}
}
