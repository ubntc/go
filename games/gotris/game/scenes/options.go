package scenes

type Options interface {
	Set(idx int)
	Get() int
	GetName() string
	List() []string
	Descs() []string
	// TODO: add SetOption(const, value)
	// to allow setting custom options according to the options
	// defined in scenes, such as scenes.OptionRenderingMode
}

type SceneOptions struct {
	Options       []string
	Descriptions  []string
	currentOption int
}

func (s *SceneOptions) List() []string {
	return s.Options
}

func (s *SceneOptions) Descs() []string {
	return s.Descriptions
}

func (s *SceneOptions) GetName() string {
	return s.Options[s.Get()]
}

func (s *SceneOptions) Set(idx int) {
	if idx >= len(s.Options) || idx < 0 {
		return
	}
	s.currentOption = idx
}

func (s *SceneOptions) Get() int {
	if len(s.Options) == 0 {
		panic("cannot access menu without options")
	}
	return s.currentOption
}
