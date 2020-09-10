package batsub_test

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/stretchr/testify/assert"
	"github.com/ubntc/go/batching/batsub"
)

type source struct {
	messages  []*pubsub.Message
	sendDelay time.Duration
	done      chan struct{}
}

// Receive is a pubsub-like Receive function to acquire messages from the source.
func (rec *source) Receive(ctx context.Context, f func(context.Context, *pubsub.Message)) error {
	go func() {
		for _, m := range rec.messages {
			f(ctx, m)
			time.Sleep(rec.sendDelay)
		}
		close(rec.done)
	}()
	<-rec.done
	return nil
}

func TestBatchedSubscription(t *testing.T) {
	type Spec struct {
		len        int           // number of test messages
		cap        int           // batch capacity
		dur        time.Duration // batch interval
		sendDelay  time.Duration // delay between test messages
		expBatches int           // number of resulting batches
		expErr     error         //
	}
	cases := map[string]Spec{
		// common cases
		"small len": {5, 2, time.Second, 0, 3, nil},
		"big len":   {1000, 10, time.Second, 0, 100, nil},
		"small cap": {100, 1, time.Second, 0, 100, nil},
		"big cap":   {100, 1000, time.Second, 0, 1, nil},
		// special cases
		"timeout":  {2, 10, time.Microsecond, time.Millisecond, 2, nil},
		"zero len": {0, 10, time.Second, 0, 0, nil},
		"zero cap": {10, 0, time.Second, 0, 10, nil},
		// NOTE: A zero capacity case is valid since Go's `append` will add missing slice capacity
		//       and the feeding of the slice will be done using a zero capacity (blocking) channel.
	}

	var wg sync.WaitGroup
	for name, spec := range cases {
		wg.Add(1)
		go func(name string, spec Spec) {
			t.Run(name, func(t *testing.T) {
				defer wg.Done()
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()

				data := make([]*pubsub.Message, 0, spec.len)
				for i := 0; i < spec.len; i++ {
					data = append(data, &pubsub.Message{ID: fmt.Sprint(i)})
				}
				assert.Len(t, data, spec.len)

				rec := &source{
					messages:  data,
					done:      make(chan struct{}),
					sendDelay: spec.sendDelay,
				}
				sub := batsub.NewBatchedSubscription(rec, spec.cap, spec.dur)

				var mu sync.Mutex
				var result []*pubsub.Message
				var numBatches = 0
				receive := func(ctx context.Context, messages []*pubsub.Message) {
					mu.Lock()
					defer mu.Unlock()
					numBatches++
					result = append(result, messages...)
				}

				err := sub.ReceiveBatches(ctx, receive)
				assert.NoError(t, err)
				assert.Len(t, result, spec.len)
				assert.Equal(t, spec.expBatches, numBatches)
			})
		}(name, spec)
	}

	wg.Wait()

}
