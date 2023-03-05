package text

import "github.com/ubntc/go/games/gotris/common/options"

type renderingOptions struct {
	// MemStore inherits the options.Options implementation from options.MemStore as embedded struct.
	options.MemStore
	p *TextUI
}

// Set overrides the Options.Set method for switching rendering
// modes when the corresponding option value changes.
func (o *renderingOptions) Set(i int) {
	o.MemStore.Set(i)
	o.p.SetRenderingMode(o.GetName())
}
