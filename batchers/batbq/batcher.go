package batbq

import (
	"context"
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
	capacity      int
	flushInterval time.Duration
	numWorkers    int
}

// NewInsertBatcher returns an InsertBatcher.
func NewInsertBatcher(capacity int, flushInterval time.Duration, numWorkers int) *InsertBatcher {
	if numWorkers <= 0 {
		numWorkers = 1
	}
	return &InsertBatcher{
		capacity:      capacity,
		flushInterval: flushInterval,
		numWorkers:    numWorkers,
	}
}

func toStructs(messages []Message) []*bigquery.StructSaver {
	res := make([]*bigquery.StructSaver, len(messages))
	for i, m := range messages {
		res[i] = m.Data()
	}
	return res
}

// Process batches messages from the given input channel to the batch-processing out Putter.
func (ins *InsertBatcher) Process(ctx context.Context, input <-chan Message, out Putter) {
	var wg sync.WaitGroup
	defer wg.Wait()

	for i := 0; i < ins.numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			Worker(ctx, ins, input, out)
		}()
	}
}
