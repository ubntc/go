package batbq

import (
	"context"
	"log"
	"sync"
	"time"
)

// Worker starts a new worker.
func Worker(ctx context.Context, ins *InsertBatcher, input <-chan Message, out Putter) {
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
		}()
	}

	ticker := time.NewTicker(ins.flushInterval)
	defer ticker.Stop()

	flush := func() {
		if len(batch) == 0 {
			return
		}
		err := out.Put(ctx, toStructs(batch))
		confirm(batch, err)                      // use current slice to ack/nack processed messages
		batch = make([]Message, 0, ins.capacity) // create a new slice to allow immediate refill
		ticker.Reset(ins.flushInterval)          // reset the ticker to avoid too early flush
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
			if len(batch) >= ins.capacity {
				flush()
			}
		}
	}
}
