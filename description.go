package metrics_adapter

import (
	"github.com/prometheus/client_golang/prometheus"
	"strings"
)

type description prometheus.Desc

// Describe sends the super-set of all possible descriptors of metrics
// collected by this Collector to the provided channel and returns once
// the last descriptor has been sent. The sent descriptors fulfill the
// consistency and uniqueness requirements described in the Desc
// documentation.
//
// It is valid if one and the same Collector sends duplicate
// descriptors. Those duplicates are simply ignored. However, two
// different Collectors must not send duplicate descriptors.
//
// Sending no descriptor at all marks the Collector as “unchecked”,
// i.e. no checks will be performed at registration time, and the
// Collector may yield any Metric it sees fit in its Collect method.
//
// This method idempotently sends the same descriptors throughout the
// lifetime of the Collector. It may be called concurrently and
// therefore must be implemented in a concurrency safe way.
//
// If a Collector encounters an error while executing this method, it
// must send an invalid descriptor (created with NewInvalidDesc) to
// signal the error to the registry.
func (d *description) Describe(ch chan<- *prometheus.Desc) {
	ch <- (*prometheus.Desc)(d)
}

// Desc returns the descriptor for the Metric. This method idempotently
// returns the same descriptor throughout the lifetime of the
// Metric. The returned descriptor is immutable by contract. A Metric
// unable to describe itself must return an invalid descriptor (created
// with NewInvalidDesc).
func (d *description) Desc() *prometheus.Desc {
	return (*prometheus.Desc)(d)
}

func newDescriptionFrom(name string) *description {
	name = strings.Replace(name, "-", "_", -1)

	return (*description)(prometheus.NewDesc(
		"witches_io_metrics_adapter_"+name,
		"Adapter for "+name,
		nil,
		nil,
	))
}
