package batsub

import (
	"context"
	"fmt"
	"sync"
	"time"

	"cloud.google.com/go/pubsub"
)

// Receiver defines a pubsub compatible `Receive` func.
type Receiver interface {
	Receive(ctx context.Context, f func(context.Context, *pubsub.Message)) error
}

// BatchedSubscription implements automatic batching
// based on a defined batch Capacity and a FlushInterval.
type BatchedSubscription struct {
	Receiver
	capacity      int
	flushInterval time.Duration
}

// NewBatchedSubscription returns an initalized BatNewBatchedSubscription.
func NewBatchedSubscription(receiver Receiver, capacity int, flushInterval time.Duration) *BatchedSubscription {
	return &BatchedSubscription{
		Receiver:      receiver,
		capacity:      capacity,
		flushInterval: flushInterval,
	}
}

// BatchFunc handles a batch of messages.
type BatchFunc func(ctx context.Context, messages []*pubsub.Message)

// ReceiveBatch calls f with the outstanding batched messages from the subscription.
//
// Basic Example:
//
// err := sub.ReceiveBatch(ctx, func(ctx context.Context, messages []*pubsub.Message){
//     for i, m := range messages {
//         // TODO: handle message
//	       m.Ack()
//     }
// })
//
// Batch Processing Example:
//
// err := sub.ReceiveBatch(ctx, func(ctx context.Context, messages []*pubsub.Message){
//
//     // handle batch of messages using a batch-processing library
//     errors := mylib.BatchProcessMessages(messages)
//     for i, err := errors {
//         if err != nil {
//             // TODO: handle error
//             continue
//         }
//         messages[i].Ack()
//     }
// })
//
func (sub *BatchedSubscription) ReceiveBatch(ctx context.Context, f BatchFunc) error {
	ch := make(chan *pubsub.Message, sub.capacity)
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		tick := time.Tick(sub.flushInterval)
		var batch []*pubsub.Message
		flush := func() {
			if len(batch) == 0 {
				return
			}
			f(ctx, batch)
			batch = make([]*pubsub.Message, 0, sub.capacity)
		}
		defer flush()

		for {
			select {
			case <-tick:
				flush()
			case msg, more := <-ch:
				if !more {
					return
				}
				batch = append(batch, msg)
				if len(batch) >= sub.capacity {
					flush()
				}
			}
		}
	}()

	// this will block until the receiver stopped
	err := sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) { ch <- msg })
	close(ch)
	wg.Wait()

	if err != nil {
		return fmt.Errorf("ReceiveBatch: %v", err)
	}

	return nil
}
