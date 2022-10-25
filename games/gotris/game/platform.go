package game

import (
	"context"

	"github.com/ubntc/go/games/gotris/game/scenes"
	"github.com/ubntc/go/games/gotris/input"
)

type Platform interface {
	// CaptureInput starts capturing input. It returns an input channel, where all input is
	// sent through, a stopCapture func to stop capturing after the game stops, and an error
	// to indicate that capturing is not possiible.
	CaptureInput(ctx context.Context) (input <-chan input.Key, stopCapture func(), err error)

	Render(game *Game)
	RenderScene(screen *scenes.Scene)
	RenderMessage(message string)

	RenderingModes() (names []string, currentMode int)
	RenderingInfo(name string) string
	SetRenderingMode(mode string) error

	// Run is a blocking call to start the platform.
	// This is the last function called to handover control to the
	// platform code. It is needed because most GUI libs need to
	// run in the main thread.
	Run(ctx context.Context)
}
