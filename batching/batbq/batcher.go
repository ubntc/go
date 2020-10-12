package batbq

import (
	"context"
	"errors"
	"sync"

	"github.com/ubntc/go/batching/batbq/config"
	"github.com/ubntc/go/batching/batbq/scaling"
)

// Putter provides a `Put` func as used by the `bigquery.Inserter`.
type Putter interface {
	Put(ctx context.Context, src interface{}) error
}

// InsertBatcher implements automatic batching with a batch capacity and flushInterval.
type InsertBatcher struct {
	id      string
	cfg     config.BatcherConfig
	metrics *Metrics
	input   <-chan Message
	output  Putter
	scaling scaling.Status
	mu      *sync.Mutex
}

// NewInsertBatcher returns an InsertBatcher.
func NewInsertBatcher(id string, opt ...BatcherOption) *InsertBatcher {
	ins := &InsertBatcher{
		id:  id,
		cfg: config.Default(),
		mu:  &sync.Mutex{},
	}
	for _, o := range opt {
		o.apply(ins)
	}
	if ins.metrics == nil {
		ins.metrics = NewMetrics()
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
		scaling.Autoscale(ctx, &ins.cfg, &ins.scaling, ins.worker)
		return nil
	}

	ins.worker(ctx, 1)
	return nil
}
