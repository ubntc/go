package batbq

import "context"

// ConfigEntry defines a batch pipeline.
type ConfigEntry struct {
	name   string
	output Putter
	input  <-chan Message
}

// Config defines multiple batch pipelines.
type Config []ConfigEntry

// MultiInsertBatcher streams data to multiple outputs.
type MultiInsertBatcher struct{}

// Process starts the batchers.
func (mb *MultiInsertBatcher) Process(ctx context.Context) {

}
