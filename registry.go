package metrics_adapter

import (
	"fmt"
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
	panic("implement me")
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
			fmt.Printf("%s %T %+v\n", name, metric, metric)
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
	if i, ok := a.names[name].(metrics.Stoppable); ok {
		i.Stop()
	}
	a.registerer.Unregister(a.names[name])
	delete(a.names, name)
}

// Unregister all metrics.  (Mostly for testing.)
func (a *registry) UnregisterAll() {
	a.mu.Lock()
	defer a.mu.Unlock()
	for name, c := range a.names {
		if i, ok := a.names[name].(metrics.Stoppable); ok {
			i.Stop()
		}
		a.registerer.Unregister(c)
		delete(a.names, name)
	}
}
