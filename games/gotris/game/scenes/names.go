package scenes

type Scene struct {
	Name          string
	Options       []string
	Descriptions  []string
	CurrentOption int
}

const (
	Welcome  = "Welcome"
	GameOver = "GameOver"
	Controls = "Controls"
	Options  = "Options"

	START, OPTIONS, CONTROLS, QUIT = "START", "OPTIONS", "CONTROLS", "QUIT"
)

func NewWelcomeMenu() *Scene {
	return &Scene{
		Name:    Welcome,
		Options: []string{START, OPTIONS, CONTROLS, QUIT},
	}
}
