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
	DefaultFlushInterval = time.Second // when to send partially filled batches
	DefaultMinWorkers    = 1

	MaxWorkerFactor = 10 // factor multiplied to cfg.MinWorkers to determine the MaxWorkers
	DrainedDivisor  = 10 // divisor applied to input channel length to check for drained channels
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
	if cfg.FlushInterval <= 0 {
		cfg.FlushInterval = DefaultFlushInterval
	}
	if cfg.MinWorkers <= 0 {
		cfg.MinWorkers = DefaultMinWorkers
	}
	if cfg.ScaleInterval <= 0 {
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
	hooks := make(map[context.Context]func())
	mu := &sync.Mutex{}

	addWorker := func() {
		mu.Lock()
		defer mu.Unlock()

		if len(hooks) >= cfg.MinWorkers*MaxWorkerFactor {
			return
		}
		log.Printf("adding worker #%d", len(hooks)+1)
		wctx, cancel := context.WithCancel(ctx)
		hooks[wctx] = cancel

		wg.Add(1)
		go func() {
			defer wg.Done()
			Worker(wctx, cfg, ins.metrics, input, out)

			mu.Lock()
			delete(hooks, wctx)
			ins.metrics.SetWorkers(len(hooks))
			mu.Unlock()
		}()

		ins.metrics.SetWorkers(len(hooks))
	}

	rmWorker := func() {
		mu.Lock()
		defer mu.Unlock()
		if len(hooks) <= cfg.MinWorkers {
			return
		}
		log.Printf("removing one of %d workers", len(hooks))
		for _, cancel := range hooks {
			cancel()
			break
		}
	}

	for i := 0; i < cfg.MinWorkers; i++ {
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
