package multibatcher

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/ubntc/go/batching/batbq"
)

// MultiBatcher streams data to multiple outputs.
type MultiBatcher struct {
	ids  []string
	opts []batbq.BatcherOption
}

// NewMultiBatcher returns a new MultiInsertBatcher
func NewMultiBatcher(ids []string, opts ...batbq.BatcherOption) *MultiBatcher {
	mib := &MultiBatcher{opts: opts}
	for _, id := range ids {
		mib.ids = append(mib.ids, id)
	}
	return mib
}

// InputGetter returns an input channel for a given batcher ID.
type InputGetter func(id string) <-chan batbq.Message

// OutputGetter returns a Putter for a given batcher ID.
type OutputGetter func(id string) batbq.Putter

// Process starts the batchers.
func (mb *MultiBatcher) Process(ctx context.Context, input InputGetter, output OutputGetter) <-chan error {
	batchers := make(map[string]*batbq.InsertBatcher)

	errchan := make(chan error, len(mb.ids))
	var wg sync.WaitGroup
	for _, id := range mb.ids {
		ins := batbq.NewInsertBatcher(id, mb.opts...)
		batchers[id] = ins
		wg.Add(1)
		go func(id string) {
			defer wg.Done()
			in := input(id)
			out := output(id)
			if err := ins.Process(ctx, in, out); err != nil {
				log.Print(err)
				errchan <- fmt.Errorf("failed to process pipeline: %s", id)
			}
		}(id)
	}

	go func() {
		wg.Wait()
		close(errchan)
	}()

	return errchan
}

// MustProcess starts the batcher and stops if any of the pipelines fails.
func (mb *MultiBatcher) MustProcess(ctx context.Context, input InputGetter, output OutputGetter) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	for err := range mb.Process(ctx, input, output) {
		return err
	}

	return nil
}
