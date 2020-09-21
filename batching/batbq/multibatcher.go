package batbq

import (
	"context"
	"fmt"
	"log"
	"sync"
)

// MultiBatcher streams data to multiple outputs.
type MultiBatcher struct {
	ids  []ID
	opts []batcherOption
}

// NewMultiBatcher returns a new MultiInsertBatcher
func NewMultiBatcher(ids []string, opts ...batcherOption) *MultiBatcher {
	mib := &MultiBatcher{opts: opts}
	for _, id := range ids {
		mib.ids = append(mib.ids, ID(id))
	}
	return mib
}

// InputGetter returns an input channel for a given batcher ID.
type InputGetter func(id ID) <-chan Message

// OutputGetter returns a Putter for a given batcher ID.
type OutputGetter func(id ID) Putter

// Process starts the batchers.
func (mb *MultiBatcher) Process(ctx context.Context, input InputGetter, output OutputGetter) <-chan error {
	batchers := make(map[ID]*InsertBatcher)

	errchan := make(chan error, len(mb.ids))
	var wg sync.WaitGroup
	for _, id := range mb.ids {
		ins := NewInsertBatcher(id, mb.opts...)
		batchers[id] = ins
		wg.Add(1)
		go func(id ID) {
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
