package batbq

import (
	"context"
	"errors"
	"sync"
	"time"

	"cloud.google.com/go/bigquery"
)

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
	workers := ins.metrics.NumWorkers.WithLabelValues(string(ins.id))
	if ins.cfg.AutoScale {
		autoscale(ctx, ins)
		return nil
	}

	workers.Inc()
	ins.worker(ctx)
	workers.Dec()
	return nil
}

func handleErrors(messages []Message, err error) map[int]struct{} {
	if err == nil {
		return nil
	}
	nacked := make(map[int]struct{})
	mulErr, isMulti := err.(bigquery.PutMultiError)
	if isMulti {
		for _, insErr := range mulErr {
			messages[insErr.RowIndex].Nack(insErr.Errors)
			nacked[insErr.RowIndex] = struct{}{}
		}
	} else {
		for i, m := range messages {
			nacked[i] = struct{}{}
			m.Nack(err)
		}
	}
	return nacked
}

func (ins *InsertBatcher) worker(ctx context.Context) {
	var wg sync.WaitGroup
	defer wg.Wait()

	var (
		batch []Message // batch of messages to be filled from the input channel

		cfg    = ins.cfg
		input  = ins.input
		output = ins.output
		name   = string(ins.id)

		insertLatency = ins.metrics.InsertLatency.WithLabelValues(name)
		ackLatency    = ins.metrics.AckLatency.WithLabelValues(name)
		errCount      = ins.metrics.InsertErrors.WithLabelValues(name)
		msgCount      = ins.metrics.ReceivedMessages.WithLabelValues(name)
		batchCount    = ins.metrics.ProcessedBatches.WithLabelValues(name)
		successCount  = ins.metrics.ProcessedMessages.WithLabelValues(name)
	)

	confirm := func(messages []Message, err error) {
		// Ensure we wait for pending (n)acks after the batcher stops.

		wg.Add(1)
		go func() {
			defer wg.Done() // allow the batcher to stop
			tStart := time.Now()
			nacked := handleErrors(messages, err)

			switch {
			case len(nacked) == len(messages):
				// all messages were nacked
			case len(nacked) == 0:
				for _, m := range messages {
					m.Ack()
				}
			default:
				for i, m := range messages {
					if _, ok := nacked[i]; ok {
						continue
					}
					m.Ack()
				}
			}
			ackLatency.Observe(time.Now().Sub(tStart).Seconds())
			successCount.Add(float64(len(messages) - len(nacked)))
			errCount.Add(float64(len(nacked)))
			batchCount.Add(1)
		}()
	}

	ticker := time.NewTicker(cfg.FlushInterval)
	defer ticker.Stop()

	flush := func() {
		if len(batch) == 0 {
			return
		}

		msgCount.Add(float64(len(batch)))
		tStart := time.Now()
		err := output.Put(ctx, toStructs(batch))
		insertLatency.Observe(time.Now().Sub(tStart).Seconds())

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
