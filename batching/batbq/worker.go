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

	flush := func() {
		if len(batch) == 0 {
			return
		}
		err := out.Put(ctx, toStructs(batch))
		batchCopy := batch                       // create a copy for async (n)acking messages
		batch = make([]Message, 0, ins.capacity) // clear the batch to allow refill
		confirm(batchCopy, err)                  // ack/nack processed messages
	}
	defer flush()

	tick := time.Tick(ins.flushInterval)

	for {
		select {
		case <-ctx.Done():
			return
		case <-tick:
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
