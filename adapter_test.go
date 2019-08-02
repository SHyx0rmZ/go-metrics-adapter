package metrics_adapter_test

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"

	metrics_adapter "code.witches.io/go/metrics-adapter"
)

func TestNilUnregister(t *testing.T) {
	r := prometheus.NewRegistry()
	a := metrics_adapter.NewRegistry(r)
	a.Unregister("any")
}
