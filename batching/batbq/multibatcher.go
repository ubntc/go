package batbq

import (
	"context"
	"fmt"
	"log"
	"sync"
)

// PipelineID defines a specific batch pipeline.
type PipelineID string

// MultiBatcher streams data to multiple outputs.
type MultiBatcher struct {
	ids []PipelineID
	cfg BatcherConfig
}

// NewMultiBatcher returns a new MultiInsertBatcher
func NewMultiBatcher(ids []string, cfg BatcherConfig) *MultiBatcher {
	mib := &MultiBatcher{
		cfg: cfg.WithDefaults(),
	}
	for _, id := range ids {
		mib.ids = append(mib.ids, PipelineID(id))
	}
	return mib
}

// InputGetter returns an input channel for a given PipelineID.
type InputGetter func(id PipelineID) <-chan Message

// OutputGetter returns a Putter for a given PipelineID.
type OutputGetter func(id PipelineID) Putter

// Process starts the batchers.
func (mb *MultiBatcher) Process(ctx context.Context, input InputGetter, output OutputGetter) <-chan error {
	batchers := make(map[PipelineID]*InsertBatcher)

	errchan := make(chan error, len(mb.ids))
	var wg sync.WaitGroup
	for _, id := range mb.ids {
		ins := NewInsertBatcher(mb.cfg)
		batchers[id] = ins
		wg.Add(1)
		go func(id PipelineID) {
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

	// TODO: go func() { collects joint metrics from `batchers` by `id` }

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
