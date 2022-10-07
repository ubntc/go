package game

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/ubntc/go/games/gotris/game/screens"
)

type Capturing bool

const (
	// capture all terminal input
	CaptureOn Capturing = true
	// do not capture terminal input (for tests)
	CaptureOff Capturing = false
)

// Run starts the main loop of the game. It also manages user input and uses
// the game's defined platform to render the game.
func (g *Game) Run(capture Capturing) error {
	// TODO: allow speedup of ticker on higher levels

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var gameOverTimeout time.Duration

	// do not capture input during tests
	if capture {
		ch, restore, err := g.platform.CaptureInput(ctx)
		if err != nil {
			return err
		}
		defer restore()
		g.input = ch
		gameOverTimeout = time.Second * 5
	}

	g.ShowScreen(screens.WelcomeScreen, 0)

	if err := g.Loop(ctx); err != nil {
		g.ShowScreen(screens.GameOver, gameOverTimeout)
	}

	// redraw the game to show the last Game state after the GameOver screen
	g.platform.Render(g)

	return nil
}

func (g *Game) Loop(ctx context.Context) error {
	var lastErr error
	ticker := time.NewTicker(g.TickTime)
	defer ticker.Stop()
	for {
		g.platform.Render(g)
		if lastErr != nil {
			g.platform.RenderMessage(lastErr.Error())
		}

		select {
		case <-ctx.Done():
			return nil
		case key, more := <-g.input:
			if !more {
				return nil
			}
			if cmd, ok := KeyToCmd(key); ok {
				lastErr = g.RunCommand(cmd)
			}
		case <-ticker.C:
			ticker.Reset(time.Duration(g.Speed))
			if err := g.Advance(); err != nil {
				lastErr = err
				return errors.Wrap(err, "GAME OVER!")
			}
			if g.MaxSteps > 0 && g.Steps > g.MaxSteps {
				return errors.New("GAME END!")
			}
		}
	}
}
