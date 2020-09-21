package batbq

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"
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

func (ins *InsertBatcher) worker(ctx context.Context, num int) {
	var wg sync.WaitGroup
	defer wg.Wait()

	var (
		batch []Message // batch of messages to be filled from the input channel

		cfg    = ins.cfg
		input  = ins.input
		output = ins.output
		name   = string(ins.id)

		workers         = ins.metrics.NumWorkers.WithLabelValues(name)
		insertLatency   = ins.metrics.InsertLatency.WithLabelValues(name)
		ackLatency      = ins.metrics.AckLatency.WithLabelValues(name)
		errCount        = ins.metrics.InsertErrors.WithLabelValues(name)
		msgCount        = ins.metrics.ReceivedMessages.WithLabelValues(name)
		batchCount      = ins.metrics.ProcessedBatches.WithLabelValues(name)
		successCount    = ins.metrics.ProcessedMessages.WithLabelValues(name)
		pendingSize     = ins.metrics.PendingMessages.WithLabelValues(name)
		pendingConfirms = ins.metrics.PendingConfirmations.WithLabelValues(name)
	)

	workers.Inc()
	defer workers.Dec()

	confirm := func(messages []Message, err error) {
		// Ensure we wait for pending (n)acks after the batcher stops.
		wg.Add(1)
		pendingConfirms.Add(float64(len(messages)))
		go func() {
			defer wg.Done() // allow the batcher to stop
			tStart := time.Now()
			acked, nacked := confirmMessages(messages, err)
			ackLatency.Observe(time.Now().Sub(tStart).Seconds())
			successCount.Add(float64(acked))
			errCount.Add(float64(nacked))
			batchCount.Add(1)
			pendingConfirms.Sub(float64(len(messages)))
		}()
	}

	ticker := time.NewTicker(cfg.FlushInterval)
	defer ticker.Stop()

	flush := func() {
		switch len(batch) {
		case 0:
			ins.scaling.Dec()
			return
		case cfg.Capacity:
			ins.scaling.Inc()
		default:
			ins.scaling.Dec()
		}

		msgCount.Add(float64(len(batch)))

		tStart := time.Now()
		err := output.Put(context.Background(), toStructs(batch))
		insertLatency.Observe(time.Now().Sub(tStart).Seconds())

		confirm(batch, err)                      // use current slice to ack/nack processed messages
		batch = make([]Message, 0, cfg.Capacity) // create a new slice to allow immediate refill
	}
	defer flush()

	log.Printf("starting worker #%d", num)
	defer log.Printf("worker #%d stopped ", num)
	for {
		pendingSize.Set(float64(len(input)))
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
