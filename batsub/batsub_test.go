package batsub

import (
	"context"
	"sync"
	"testing"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/stretchr/testify/assert"
)

type testReceiver struct {
	messages []*pubsub.Message
	done     chan struct{}
}

func (rec *testReceiver) Receive(ctx context.Context, f func(context.Context, *pubsub.Message)) error {
	go func() {
		for _, m := range rec.messages {
			f(ctx, m)
		}
		close(rec.done)
	}()
	<-ctx.Done()
	return nil
}

func TestBatchedSubscription(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	rec := &testReceiver{
		messages: []*pubsub.Message{
			{ID: "1", Data: []byte("a")},
			{ID: "2", Data: []byte("b")},
			{ID: "3", Data: []byte("c")},
			{ID: "4", Data: []byte("d")},
			{ID: "5", Data: []byte("e")},
		},
		done: make(chan struct{}),
	}
	sub := NewBatchedSubscription(rec, 2, time.Millisecond)

	var mu sync.Mutex
	var result []*pubsub.Message
	var numBatches = 0
	receive := func(ctx context.Context, messages []*pubsub.Message) {
		mu.Lock()
		defer mu.Unlock()
		numBatches++
		result = append(result, messages...)
	}

	go func() {
		<-rec.done
		cancel()
		<-ctx.Done()
	}()

	err := sub.ReceiveBatch(ctx, receive)
	assert.NoError(t, err)
	assert.Len(t, result, 5)
	assert.Equal(t, 3, numBatches)
}
