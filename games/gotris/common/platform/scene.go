package platform

import "github.com/ubntc/go/games/gotris/common/options"

type Scene interface {
	Options() options.Options
	Name() string
}

type scene struct {
	name    string
	options options.Options
}

// NewScene returns a named Scene without any Options.
func NewScene(name string) *scene {
	return &scene{
		name: name,
	}
}

// New returns a named Scene without any Options.
func NewMenu(name string, opt options.Options) *scene {
	return &scene{
		name:    name,
		options: opt,
	}
}

func (s *scene) Options() options.Options {
	return s.options
}

func (s *scene) Name() string {
	return s.name
}
