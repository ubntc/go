package scenes

import "github.com/ubntc/go/games/gotris/game/options"

type Scene struct {
	name    string
	options options.Options
}

const (
	TitleWelcome  = "Welcome"
	TitleGameOver = "GameOver"
	TitleControls = "Controls"
	TitleOptions  = "Options"

	OptionRenderingMode = "Rendering Mode"

	START, OPTIONS, CONTROLS, QUIT = "START", "OPTIONS", "CONTROLS", "QUIT"
)

// New returns a named Scene without any Options.
func New(name string) *Scene {
	return &Scene{
		name: name,
	}
}

// New returns a named Scene without any Options.
func NewMenu(name string, opt options.Options) *Scene {
	return &Scene{
		name:    name,
		options: opt,
	}
}

func (s *Scene) Options() options.Options {
	return s.options
}

func (s *Scene) Name() string {
	return s.name
}
