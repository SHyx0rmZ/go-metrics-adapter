package metrics_adapter

import (
	"github.com/rcrowley/go-metrics"
	"sync"
	"time"
	_ "unsafe"
)

type StandardMeter metrics.StandardMeter

//go:linkname arbiter github.com/rcrowley/go-metrics.arbiter
var arbiter struct {
	sync.RWMutex
	started bool
	meters  map[*StandardMeter]struct{}
	ticker  *time.Ticker
}
