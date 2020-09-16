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

	if ins.cfg.AutoScale {
		ins.autoscale(ctx)
		return nil
	}

	ins.worker(ctx)
	return nil
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

		workers       = ins.metrics.NumWorkers.WithLabelValues(name)
		insertLatency = ins.metrics.InsertLatency.WithLabelValues(name)
		ackLatency    = ins.metrics.AckLatency.WithLabelValues(name)
		errCount      = ins.metrics.InsertErrors.WithLabelValues(name)
		msgCount      = ins.metrics.ReceivedMessages.WithLabelValues(name)
		batchCount    = ins.metrics.ProcessedBatches.WithLabelValues(name)
		successCount  = ins.metrics.ProcessedMessages.WithLabelValues(name)
	)

	workers.Inc()
	defer workers.Dec()

	confirm := func(messages []Message, err error) {
		// Ensure we wait for pending (n)acks after the batcher stops.
		wg.Add(1)
		go func() {
			defer wg.Done() // allow the batcher to stop
			tStart := time.Now()
			acked, nacked := confirmMessages(messages, err)
			ackLatency.Observe(time.Now().Sub(tStart).Seconds())
			successCount.Add(float64(acked))
			errCount.Add(float64(nacked))
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
		ticker.Reset(cfg.FlushInterval)          // reset the ticker to avoid too early next flush
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

// confirmMessages acks and nacks `messages` in the context of a potential
// batching `error` and returns the number of acked and nacked messages.
func confirmMessages(messages []Message, err error) (numAcked int, numNacked int) {
	nacked := handleErrors(messages, err)

	switch {
	case len(nacked) == len(messages):
		// all messages had errors were nacked
	case len(nacked) == 0:
		// no messages had errors and can be acked
		for _, m := range messages {
			m.Ack()
		}
	default:
		// some messages had errors, we need to check which
		for i, m := range messages {
			if _, ok := nacked[i]; ok {
				continue
			}
			m.Ack()
		}
	}
	return len(messages) - len(nacked), len(nacked)
}

// handleErrors nacks `messages` according to the type of the received `error`.
// It returns an index of the nacked messages.
func handleErrors(messages []Message, err error) (index map[int]struct{}) {
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
