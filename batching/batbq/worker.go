package batbq

import (
	"context"
	"log"
	"sync"
	"time"

	"cloud.google.com/go/bigquery"
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

		workers       = ins.metrics.NumWorkers.WithLabelValues(name)
		insertLatency = ins.metrics.InsertLatency.WithLabelValues(name)
		ackLatency    = ins.metrics.AckLatency.WithLabelValues(name)
		errCount      = ins.metrics.InsertErrors.WithLabelValues(name)
		msgCount      = ins.metrics.ReceivedMessages.WithLabelValues(name)
		batchCount    = ins.metrics.ProcessedBatches.WithLabelValues(name)
		successCount  = ins.metrics.ProcessedMessages.WithLabelValues(name)
		pendingSize   = ins.metrics.PendingMessages.WithLabelValues(name)
	)

	workers.Inc()
	defer workers.Dec()

	confirm := func(messages []Message, err error) {
		tStart := time.Now()

		acked, nacked := confirmMessages(messages, err)

		ackLatency.Observe(time.Now().Sub(tStart).Seconds())
		successCount.Add(float64(acked))
		errCount.Add(float64(nacked))
		batchCount.Add(1)
	}

	put := func(messages []Message) error {
		tStart := time.Now()

		rows := make([]bigquery.ValueSaver, len(messages))
		for i, m := range messages {
			rows[i] = m.Data()
		}
		err := output.Put(context.Background(), rows)

		insertLatency.Observe(time.Now().Sub(tStart).Seconds())
		return err
	}

	flush := func() {
		if cfg.AutoScale {
			ins.scaling.UpdateLoadLevel(len(batch), len(input), cfg.Capacity)
		}

		if len(batch) == 0 {
			return
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

	ticker := time.NewTicker(cfg.FlushInterval)
	defer ticker.Stop()

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
				// Make sure we do not flush more than once per second sending unfilled batches
				// since this will confuse the autoscaler.
				ticker.Reset(cfg.FlushInterval)
			}
		}
	}
}
