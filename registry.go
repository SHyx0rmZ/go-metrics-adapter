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
	names      map[string]prometheus.Collector
	mu         sync.Mutex
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
	for s, c := range a.names {
		f(s, c)
	}
}

func (a *registry) Get(s string) interface{} {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.names[s]
}

func (a *registry) GetAll() map[string]map[string]interface{} {
	panic("implement me")
}

func (a *registry) GetOrRegister(s string, v interface{}) interface{} {
	c := a.Get(s)
	if c == nil {
		err := a.Register(s, v)
		if err != nil {
			panic(err)
		}
		return a.Get(s)
	}
	return c
}

func (a *registry) Register(s string, v interface{}) error {
	c, ok := v.(prometheus.Collector)
	if !ok {
		if reflect.TypeOf(v).Kind() == reflect.Func {
			v = reflect.ValueOf(v).Call(nil)[0].Interface()
		}
		switch v := v.(type) {
		case metrics.Counter:
			c = NewCounter(s, v)
		case metrics.Gauge:
			c = NewGauge(s, v)
		case metrics.GaugeFloat64:
			c = NewGaugeFloat64(s, v)
		case metrics.Healthcheck:
			c = NewHealthcheck(s, v)
		case metrics.Histogram:
			c = NewHistogram(s, v)
		case metrics.Meter:
			c = NewMeter(s, v)
		case metrics.Timer:
			c = NewTimer(s, v)
		default:
			fmt.Printf("%s %T %+v\n", s, v, v)
			return ErrExpectedCollector
		}
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	err := a.registerer.Register(c)
	if err != nil {
		return err
	}
	a.names[s] = c
	return nil
}

func (a *registry) RunHealthchecks() {
	a.Each(func(s string, i interface{}) {
		if h, ok := i.(metrics.Healthcheck); ok {
			h.Check()
		}
	})
}

func (a *registry) Unregister(s string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if i, ok := a.names[s].(metrics.Stoppable); ok {
		i.Stop()
	}
	a.registerer.Unregister(a.names[s])
	delete(a.names, s)
}

func (a *registry) UnregisterAll() {
	a.mu.Lock()
	defer a.mu.Unlock()
	for s, c := range a.names {
		if i, ok := a.names[s].(metrics.Stoppable); ok {
			i.Stop()
		}
		a.registerer.Unregister(c)
		delete(a.names, s)
	}
}
