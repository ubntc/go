package batbq

import (
	"context"
	"log"
	"sync"
	"time"
)

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
		pendingConfirms.Add(float64(len(messages)))

		tStart := time.Now()
		acked, nacked := confirmMessages(messages, err)
		ackLatency.Observe(time.Now().Sub(tStart).Seconds())

		successCount.Add(float64(acked))
		errCount.Add(float64(nacked))
		batchCount.Add(1)
		pendingConfirms.Sub(float64(len(messages)))
	}

	put := func(messages []Message) error {
		tStart := time.Now()
		err := output.Put(context.Background(), toStructs(messages))
		insertLatency.Observe(time.Now().Sub(tStart).Seconds())
		return err
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

		wg.Add(1) // Ensure we wait for pending puts and (n)acks after the batcher stops.
		go func(messages []Message) {
			defer wg.Done() // Allow the batcher to stop after the last batch was processed.
			err := put(messages)
			confirm(messages, err)
		}(batch)

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
