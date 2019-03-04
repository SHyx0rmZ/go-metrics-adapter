package metrics_adapter

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rcrowley/go-metrics"
	"testing"
)

func BenchmarkRegistry(b *testing.B) {
	r := NewRegistry(prometheus.NewRegistry())
	r.Register("foo", metrics.NewCounter())
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r.Each(func(string, interface{}) {})
	}
}

func BenchmarkRegistryParallel(b *testing.B) {
	r := NewRegistry(prometheus.NewRegistry())
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			r.GetOrRegister("foo", metrics.NewCounter())
		}
	})
}

func TestRegistry(t *testing.T) {
	r := NewRegistry(prometheus.NewRegistry())
	r.Register("foo", metrics.NewCounter())
	i := 0
	r.Each(func(name string, iface interface{}) {
		i++
		if "foo" != name {
			t.Fatal(name)
		}
		if _, ok := iface.(metrics.Counter); !ok {
			t.Fatal(iface)
		}
	})
	if 1 != i {
		t.Fatal(i)
	}
	r.Unregister("foo")
	i = 0
	r.Each(func(string, interface{}) { i++ })
	if 0 != i {
		t.Fatal(i)
	}
}

func TestRegistryDuplicate(t *testing.T) {
	r := NewRegistry(prometheus.NewRegistry())
	if err := r.Register("foo", metrics.NewCounter()); nil != err {
		t.Fatal(err)
	}
	if err := r.Register("foo", metrics.NewGauge()); nil == err {
		t.Fatal(err)
	}
	i := 0
	r.Each(func(name string, iface interface{}) {
		i++
		if _, ok := iface.(metrics.Counter); !ok {
			t.Fatal(iface)
		}
	})
	if 1 != i {
		t.Fatal(i)
	}
}

func TestRegistryGet(t *testing.T) {
	r := NewRegistry(prometheus.NewRegistry())
	r.Register("foo", metrics.NewCounter())
	if count := r.Get("foo").(metrics.Counter).Count(); 0 != count {
		t.Fatal(count)
	}
	r.Get("foo").(metrics.Counter).Inc(1)
	if count := r.Get("foo").(metrics.Counter).Count(); 1 != count {
		t.Fatal(count)
	}
}

func TestRegistryGetOrRegister(t *testing.T) {
	r := NewRegistry(prometheus.NewRegistry())

	// First metric wins with GetOrRegister
	_ = r.GetOrRegister("foo", metrics.NewCounter())
	m := r.GetOrRegister("foo", metrics.NewGauge())
	if _, ok := m.(metrics.Counter); !ok {
		t.Fatal(m)
	}

	i := 0
	r.Each(func(name string, iface interface{}) {
		i++
		if name != "foo" {
			t.Fatal(name)
		}
		if _, ok := iface.(metrics.Counter); !ok {
			t.Fatal(iface)
		}
	})
	if i != 1 {
		t.Fatal(i)
	}
}

func TestRegistryGetOrRegisterWithLazyInstantiation(t *testing.T) {
	r := NewRegistry(prometheus.NewRegistry())

	// First metric wins with GetOrRegister
	_ = r.GetOrRegister("foo", metrics.NewCounter)
	m := r.GetOrRegister("foo", metrics.NewGauge)
	if _, ok := m.(metrics.Counter); !ok {
		t.Fatal(m)
	}

	i := 0
	r.Each(func(name string, iface interface{}) {
		i++
		if name != "foo" {
			t.Fatal(name)
		}
		if _, ok := iface.(metrics.Counter); !ok {
			t.Fatal(iface)
		}
	})
	if i != 1 {
		t.Fatal(i)
	}
}

func TestRegistryUnregister(t *testing.T) {
	l := len(arbiter.meters)
	r := NewRegistry(prometheus.NewRegistry())
	r.Register("foo", metrics.NewCounter())
	r.Register("bar", metrics.NewMeter())
	r.Register("baz", metrics.NewTimer())
	if len(arbiter.meters) != l+2 {
		t.Errorf("arbiter.meters: %d != %d\n", l+2, len(arbiter.meters))
	}
	r.Unregister("foo")
	r.Unregister("bar")
	r.Unregister("baz")
	if len(arbiter.meters) != l {
		t.Errorf("arbiter.meters: %d != %d\n", l+2, len(arbiter.meters))
	}
}
