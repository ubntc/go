package game

import (
	"context"

	"github.com/ubntc/go/games/gotris/input"
)

type Platform interface {
	CaptureInput(ctx context.Context) (<-chan input.Key, func(), error)

	Render(game *Game)
	RenderScreen(textArt string)
	RenderMessage(message string)

	RenderingModes() (names []string, currentMode int)
	RenderingInfo(name string) string
	SetRenderingMode(mode string) error
}
