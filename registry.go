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

func NewRegistry(registerer prometheus.Registerer) metrics.Registry {
	return &registry{
		registerer: registerer,
		names:      make(map[string]prometheus.Collector),
	}
}

func (a *registry) Each(f func(string, interface{})) {
	a.mu.Lock()
	defer a.mu.Unlock()
	for name, collector := range a.names {
		f(name, collector)
	}
}

func (a *registry) Get(name string) interface{} {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.names[name]
}

func (a *registry) GetAll() map[string]map[string]interface{} {
	panic("implement me")
}

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

func (a *registry) Register(name string, v interface{}) error {
	c, ok := v.(prometheus.Collector)
	if !ok {
		if reflect.TypeOf(v).Kind() == reflect.Func {
			v = reflect.ValueOf(v).Call(nil)[0].Interface()
		}
		switch v := v.(type) {
		case metrics.Counter:
			c = NewCounter(name, v)
		case metrics.Gauge:
			c = NewGauge(name, v)
		case metrics.GaugeFloat64:
			c = NewGaugeFloat64(name, v)
		case metrics.Healthcheck:
			c = NewHealthcheck(name, v)
		case metrics.Histogram:
			c = NewHistogram(name, v)
		case metrics.Meter:
			c = NewMeter(name, v)
		case metrics.Timer:
			c = NewTimer(name, v)
		default:
			fmt.Printf("%s %T %+v\n", name, v, v)
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

func (a *registry) RunHealthchecks() {
	a.Each(func(name string, v interface{}) {
		if h, ok := v.(metrics.Healthcheck); ok {
			h.Check()
		}
	})
}

func (a *registry) Unregister(name string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if i, ok := a.names[name].(metrics.Stoppable); ok {
		i.Stop()
	}
	a.registerer.Unregister(a.names[name])
	delete(a.names, name)
}

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
