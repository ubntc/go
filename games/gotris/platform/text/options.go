package text

import "github.com/ubntc/go/games/gotris/game/options"

type RenderingOptions struct {
	// MemStore inherits the options.Options implementation from options.MemStore as embedded struct.
	options.MemStore
	p *Platform
}

// Set overrides the SceneOptions.Set to allow switching rendering
// modes when the option changes.
func (o *RenderingOptions) Set(i int) {
	o.MemStore.Set(i)
	o.p.SetRenderingMode(o.GetName())
}
