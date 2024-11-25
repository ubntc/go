package batsub

import (
	"context"
	"sync"
	"time"

	"cloud.google.com/go/pubsub"
)

// DefaultFlushInterval defines the default time between forced flushes of incomplete batches.
var DefaultFlushInterval = time.Second

// Receiver defines pubsub.Subscription compatible interface with a `Receive` and `ID` method.
type Receiver interface {
	// Receive handles receiving messages from a subscription.
	Receive(ctx context.Context, f func(context.Context, *pubsub.Message)) error
	// ID returns an identifier for a subscription and is used for the metrics.
	ID() string
}

// BatchedSubscription implements automatic batching
// based on a defined batch Capacity and a FlushInterval.
type BatchedSubscription struct {
	Receiver
	capacity      int
	flushInterval time.Duration
	metrics       *Metrics
}

// NewBatchedSubscription returns an initalized BatchedSubscription.
func NewBatchedSubscription(receiver Receiver, opt ...Option) *BatchedSubscription {
	b := &BatchedSubscription{Receiver: receiver}
	for _, o := range opt {
		o.apply(b)
	}
	if b.metrics == nil {
		b.metrics = NewMetrics()
	}
	if b.flushInterval == 0 {
		b.flushInterval = DefaultFlushInterval
	}
	return b
}

// BatchFunc handles a batch of messages.
type BatchFunc func(ctx context.Context, messages []*pubsub.Message)

// ReceiveBatches calls f with the outstanding batched messages from the subscription.
//
// Basic Example:
//
//	err := sub.ReceiveBatches(ctx, func(ctx context.Context, messages []*pubsub.Message){
//	    for i, m := range messages {
//	        // TODO: handle message
//		       m.Ack()
//	    }
//	})
//
// Batch Processing Example:
//
// err := sub.ReceiveBatches(ctx, func(ctx context.Context, messages []*pubsub.Message){
//
//	    // handle batch of messages using a batch-processing library
//	    errors := mylib.BatchProcessMessages(messages)
//	    for i, err := errors {
//	        if err != nil {
//	            // TODO: handle error
//	            continue
//	        }
//	        messages[i].Ack()
//	    }
//	})
func (sub *BatchedSubscription) ReceiveBatches(ctx context.Context, f BatchFunc) error {
	ch := make(chan *pubsub.Message, sub.capacity)
	var wg sync.WaitGroup

	var (
		id        = sub.ID()
		pending   = sub.metrics.PendingMessages.WithLabelValues(id)
		processed = sub.metrics.ProcessedMessages.WithLabelValues(id)
		flushed   = sub.metrics.ProcessedBatches.WithLabelValues(id)
		latency   = sub.metrics.ProcessingLatency.WithLabelValues(id)
	)

	wg.Add(1)
	go func() {
		defer wg.Done()

		var batch []*pubsub.Message
		flush := func() {
			if len(batch) == 0 {
				return
			}
			wg.Add(1) // ensure we wait for pending flushes
			go func(batch []*pubsub.Message) {
				defer wg.Done()
				f(ctx, batch)
				// track progress after batch is completed
				latency.Observe(float64(time.Since(time.Now())))
				flushed.Add(float64(len(batch)))
				processed.Add(float64(len(batch)))
			}(batch)

			batch = make([]*pubsub.Message, 0, sub.capacity)
		}
		defer flush()

		tick := time.NewTicker(sub.flushInterval)
		defer tick.Stop()
		for {
			select {
			case <-tick.C:
				flush()
			case msg, more := <-ch:
				// The batching can be stopped by closing the channel.
				// No additional context handling is required.
				if !more {
					return
				}
				batch = append(batch, msg)
				if len(batch) >= sub.capacity {
					flush()
					pending.Set(float64(len(ch)))
				}
			}
		}
	}()

	// The receiver will block until it is stopped via the external context.
	// After it stopped, no more messages must be sent to the channel.
	err := sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) { ch <- msg })

	// The channel can now be closed safely to stop the batching goroutine.
	close(ch)

	// Wait until all pending messages are flushed and all pending flushes are completed.
	wg.Wait()

	if err != nil {
		return err
	}

	return nil
}
