package scenes

type Scene interface {
	Name() string
	Options() Options
}

type SimpleScene struct {
	SceneOptions Options
	name         string
}

const (
	TitleWelcome  = "Welcome"
	TitleGameOver = "GameOver"
	TitleControls = "Controls"
	TitleOptions  = "Options"

	OptionRenderingMode = "Rendering Mode"

	START, OPTIONS, CONTROLS, QUIT = "START", "OPTIONS", "CONTROLS", "QUIT"
)

func New(name string) *SimpleScene {
	return &SimpleScene{
		name: name,
	}
}

func NewOptionsScene(name string, opt Options) *SimpleScene {
	return &SimpleScene{
		name:         name,
		SceneOptions: opt,
	}
}

func (s *SimpleScene) Options() Options {
	return s.SceneOptions
}

func (s *SimpleScene) Name() string {
	return s.name
}
