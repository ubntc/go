package batbq

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"
)

// Putter provides a `Put` func as used by the `bigquery.Inserter`.
type Putter interface {
	Put(ctx context.Context, src interface{}) error
}

// InsertBatcher implements automatic batching with a batch capacity and flushInterval.
type InsertBatcher struct {
	cfg     BatcherConfig
	metrics *metricsRecorder
	input   <-chan Message
	output  Putter
	mu      *sync.Mutex
}

// NewInsertBatcher returns an InsertBatcher.
func NewInsertBatcher(cfg BatcherConfig) *InsertBatcher {
	return &InsertBatcher{
		cfg:     cfg.WithDefaults(),
		metrics: newMetricsRecorder(),
		mu:      &sync.Mutex{},
	}
}

// Metrics returns a copy of the metrics.
func (ins *InsertBatcher) Metrics() *Metrics {
	return ins.metrics.Metrics()
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
		autoscale(ctx, ins)
		return nil
	}
	ins.metrics.SetWorkers(1)
	ins.worker(ctx)
	ins.metrics.SetWorkers(0)
	return nil
}

func (ins *InsertBatcher) worker(ctx context.Context) {
	var wg sync.WaitGroup
	defer wg.Wait()

	// batch of messages to be filled from the input channel
	var batch []Message
	cfg := ins.cfg
	input := ins.input
	output := ins.output

	var ackLock sync.Mutex
	confirm := func(messages []Message, err error) {
		ackLock.Lock()
		// The lock ensures to stop processing if the previous acks are not complete.
		// Otherwise the batcher could eat up the memory with pending unacked messages
		// in case the (n)acking takes too long.

		// Also ensure we wait for pending (n)acks after the batcher stops.
		wg.Add(1)
		go func() {
			defer wg.Done()        // allow the batcher to stop
			defer ackLock.Unlock() // allow confirming the next batch
			// TODO: handle insert errors
			if err != nil {
				log.Print(err)
			}
			for _, m := range messages {
				m.Ack()
			}
			ins.metrics.IncMessages(len(messages), 1)
		}()
	}

	ticker := time.NewTicker(cfg.FlushInterval)
	defer ticker.Stop()

	flush := func() {
		if len(batch) == 0 {
			return
		}
		err := output.Put(ctx, toStructs(batch))
		confirm(batch, err)                      // use current slice to ack/nack processed messages
		batch = make([]Message, 0, cfg.Capacity) // create a new slice to allow immediate refill
		ticker.Reset(cfg.FlushInterval)          // reset the ticker to avoid too early flush
	}
	defer flush()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			flush()
		case msg, more := <-input:
			if !more {
				return
			}
			batch = append(batch, msg)
			if len(batch) >= cfg.Capacity {
				flush()
			}
		}
	}
}
