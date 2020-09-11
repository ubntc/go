package batbq

import (
	"context"
	"log"
	"sync"
	"time"
)

// BatcherConfig defaults.
const (
	DefaultScaleInterval = time.Second // how often to trigger worker scaling
	DefaultMinWorkers    = 1
	MaxWorkerFactor      = 10 // factor multiplied to cfg.MinWorkers to determine the MaxWorkers
	DrainedDivisor       = 10 // divisor applied to input channel length to check for drained channels
)

// Putter provides a `Put` func as used by the `bigquery.Inserter`.
type Putter interface {
	Put(ctx context.Context, src interface{}) error
}

// BatcherConfig stores InsertBatcher paramaters.
type BatcherConfig struct {
	Capacity      int
	FlushInterval time.Duration
	MinWorkers    int
	ScaleInterval time.Duration
}

// InsertBatcher implements automatic batching with a batch capacity and flushInterval.
type InsertBatcher struct {
	cfg     BatcherConfig
	metrics *metricsRecorder
}

// NewInsertBatcher returns an InsertBatcher.
func NewInsertBatcher(cfg BatcherConfig) *InsertBatcher {
	if cfg.MinWorkers <= 0 {
		cfg.MinWorkers = DefaultMinWorkers
	}
	if cfg.ScaleInterval == 0 {
		cfg.ScaleInterval = DefaultScaleInterval
	}
	return &InsertBatcher{
		cfg:     cfg,
		metrics: newMetricsRecorder(),
	}
}

// Metrics returns a copy of the metrics.
func (ins *InsertBatcher) Metrics() *Metrics {
	return ins.metrics.Metrics()
}

// Process batches messages from the given input channel to the batch-processing out Putter.
func (ins *InsertBatcher) Process(ctx context.Context, input <-chan Message, out Putter) {
	var wg sync.WaitGroup

	cfg := ins.cfg

	var hooks []func()
	addWorker := func() {
		log.Printf("adding worker #%d", len(hooks)+1)
		wg.Add(1)
		wctx, cancel := context.WithCancel(ctx)
		hooks = append(hooks, cancel)
		go func() {
			defer wg.Done()
			Worker(wctx, cfg, ins.metrics, input, out)
		}()
		ins.metrics.SetWorkers(len(hooks))
	}

	rmWorker := func() {
		if len(hooks) <= cfg.MinWorkers {
			return
		}
		if len(hooks) >= cfg.MinWorkers*MaxWorkerFactor {
			return
		}
		log.Printf("removing first worker of %d workers", len(hooks))
		cancel := hooks[0]
		hooks = hooks[1:]
		cancel()
		ins.metrics.SetWorkers(len(hooks))
	}

	for len(hooks) < cfg.MinWorkers {
		addWorker()
	}

	// start worker scaling
	tick := time.NewTicker(cfg.ScaleInterval)

	go func() {
		if cfg.Capacity <= 1 {
			// cannot do capacity-based scaling for small capacities
			return
		}
		for {
			<-tick.C
			switch {
			case len(input) >= cfg.Capacity:
				addWorker()
			case len(input) < cfg.Capacity/DrainedDivisor:
				rmWorker()
			}
		}
	}()

	wg.Wait()   // wait for all workers to finish
	tick.Stop() // stop worker scaling
}
