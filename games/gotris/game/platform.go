package game

import (
	"context"

	"github.com/ubntc/go/games/gotris/input"
)

type Platform interface {
	CaptureInput(context.Context) (<-chan input.Key, func(), error)
	Render(*Game)
	RenderScreen(string)
	RenderMessage(string)
}
