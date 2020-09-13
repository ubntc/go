package batbq

import (
	"context"
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
func (mb *MultiBatcher) Process(ctx context.Context, input InputGetter, output OutputGetter) {
	var wg sync.WaitGroup
	defer wg.Wait()
	batchers := make(map[PipelineID]*InsertBatcher)
	for _, id := range mb.ids {
		ins := NewInsertBatcher(mb.cfg)
		batchers[id] = ins
		wg.Add(1)
		go func(id PipelineID) {
			defer wg.Done()
			in := input(id)
			if in == nil {
				log.Println("failed to create input channel for PipelineID:", id)
				return
			}
			out := output(id)
			if out == nil {
				log.Println("failed to create output for PipelineID:", id)
				return
			}
			ins.Process(ctx, in, out)
		}(id)
	}
	// TODO: go func() { collects joint metrics }
}
