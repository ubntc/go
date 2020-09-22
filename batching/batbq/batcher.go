package batbq

import (
	"context"
	"errors"
	"sync"
)

// ID defines a specific batch pipeline.
type ID string

// Putter provides a `Put` func as used by the `bigquery.Inserter`.
type Putter interface {
	Put(ctx context.Context, src interface{}) error
}

// InsertBatcher implements automatic batching with a batch capacity and flushInterval.
type InsertBatcher struct {
	id      ID
	cfg     BatcherConfig
	metrics *Metrics
	input   <-chan Message
	output  Putter
	scaling scalingStatus
	mu      *sync.Mutex
}

type batcherOption interface {
	Apply(*InsertBatcher)
}

// NewInsertBatcher returns an InsertBatcher.
func NewInsertBatcher(id ID, opt ...batcherOption) *InsertBatcher {
	ins := &InsertBatcher{
		id:  id,
		cfg: BatcherConfig{}.WithDefaults(),
		mu:  &sync.Mutex{},
	}
	for _, o := range opt {
		o.Apply(ins)
	}
	if ins.metrics == nil {
		ins.metrics = newMetrics()
	}
	return ins
}

// Metrics returns the metrics.
func (ins *InsertBatcher) Metrics() *Metrics {
	return ins.metrics
}

// Process starts the batcher.
func (ins *InsertBatcher) Process(ctx context.Context, input <-chan Message, output Putter) error {
	if input == nil {
		return errors.New("input channel must not be nil")
	}
	if output == nil {
		return errors.New("output Putter must not be nil")
	}
	// ensure ins.Process is not called concurrently
	ins.mu.Lock()
	defer ins.mu.Unlock()
	ins.input = input
	ins.output = output

	if ins.cfg.AutoScale {
		ins.autoscale(ctx)
		return nil
	}

	ins.worker(ctx, 1)
	return nil
}
