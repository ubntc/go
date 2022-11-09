package scenes

import "github.com/ubntc/go/games/gotris/game/options"

func NewWelcomeMenu() *Scene {
	return NewMenu(
		TitleWelcome,
		options.NewMemStore([]string{START, OPTIONS, CONTROLS, QUIT}, nil),
	)
}
