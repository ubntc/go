package scenes

func NewWelcomeMenu() Scene {
	return NewOptionsScene(
		TitleWelcome,
		&SceneOptions{
			Options: []string{START, OPTIONS, CONTROLS, QUIT},
		},
	)
}
