package game

import (
	"context"

	"github.com/ubntc/go/games/gotris/game/options"
	"github.com/ubntc/go/games/gotris/game/scenes"
	"github.com/ubntc/go/games/gotris/input"
)

type Platform interface {
	// CaptureInput starts capturing input. It returns an input channel, where all input is
	// sent through, a stopCapture func to stop capturing after the game stops, and an error
	// to indicate that capturing is not possiible.
	CaptureInput(ctx context.Context) (input <-chan *input.Input, stopCapture func(), err error)

	// ShowMessage shows a generic info message to the user.
	ShowMessage(message string)

	// Render is called by the game when the game state updates. An implementation of this
	// should update it's own presentation state based on the changes inside the game.
	Render(game *Game)

	// RenderScene is called by the game to render an specific gamer Scene, with specific
	// options, descriptions, and a given current selection. There is not feedback channel.
	// The game must be informed via regular input events, which the game will wait for after
	// calling RenderScene.
	RenderScene(scene *scenes.Scene)

	// Options returns a scene with options that can be managed by the platform.
	Options() options.Options

	// Run is a blocking call to start the platform.
	// This is the last function called to handover control to the
	// platform code. It is needed because most GUI libs need to
	// run in the main thread.
	Run(ctx context.Context)
}
