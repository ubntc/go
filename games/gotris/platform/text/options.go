package text

import "github.com/ubntc/go/games/gotris/game/scenes"

type RenderingOptions struct {
	scenes.SceneOptions
	p *Platform
}

// Set overrides the SceneOptions.Set to allow switching rendering
// modes when the option changes.
func (o *RenderingOptions) Set(i int) {
	o.SceneOptions.Set(i)
	o.p.SetRenderingMode(o.GetName())
}
