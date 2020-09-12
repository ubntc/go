package batbq

import (
	"context"
	"log"
	"sync"
	"time"

	"cloud.google.com/go/bigquery"
)

func toStructs(messages []Message) []*bigquery.StructSaver {
	res := make([]*bigquery.StructSaver, len(messages))
	for i, m := range messages {
		res[i] = m.Data()
	}
	return res
}

// Worker starts a new worker.
func Worker(ctx context.Context, cfg BatcherConfig, metrics *metricsRecorder, input <-chan Message, out Putter) {
	var wg sync.WaitGroup
	defer wg.Wait()

	// batch of messages to be filled from the input channel
	var batch []Message

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
			metrics.SetMessages(len(messages), 1)
		}()
	}

	ticker := time.NewTicker(cfg.FlushInterval)
	defer ticker.Stop()

	flush := func() {
		if len(batch) == 0 {
			return
		}
		err := out.Put(ctx, toStructs(batch))
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
