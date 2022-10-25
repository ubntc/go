package game

import (
	"context"
	"time"
)

func (g *Game) startCapture(ctx context.Context) context.CancelFunc {
	// do not capture input during tests
	if g.CaptureInput {
		ch, restore, err := g.platform.CaptureInput(ctx)
		if err != nil {
			panic(err)
		}
		g.input = ch
		g.GameOverScreenDuration = time.Second * 5
		return restore
	}
	return func() {}
}

func (g *Game) Run(ctx context.Context) {
	stopCapture := g.startCapture(ctx)
	defer stopCapture()

	g.showWelcomeScreen(ctx)
}
