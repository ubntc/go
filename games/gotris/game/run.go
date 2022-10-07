package game

import (
	"context"
	"log"
	"time"

	"github.com/pkg/errors"
	"github.com/ubntc/go/games/gotris/game/screens"
)

type Capturing bool

const (
	// capture all terminal input
	CaptureOn Capturing = true
	// do not capture all terminal input (for tests)
	CaptureOff Capturing = false
)

type Platform interface {
	CaptureInput(context.Context) (<-chan []rune, func(), error)
	Render(*Game)
	RenderText(string)
}

// Run is the main Game Run, managing user input
// and state changes in the game in a step by step way.
func Run(rules Rules, platform Platform, capture Capturing) {
	// TODO: allow speedup of ticker on higher levels
	g := NewGame(rules)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var runes <-chan []rune

	// do not capture input during tests
	if capture {
		ch, restore, err := platform.CaptureInput(ctx)
		if err != nil {
			log.Fatalln(err)
		}
		defer restore()
		runes = ch
	}

	showScreen := func(text string) {
		platform.RenderText(text)
		if capture {
			<-runes
		}
	}

	ticker := time.NewTicker(g.TickTime)
	defer ticker.Stop()

	var keys []rune
	var more bool
	var lastError error

	showScreen(screens.WelcomeScreen)

	for {
		platform.Render(g)
		g.Message = map[string]interface{}{
			// "score": g.Score,
			// "keys":  keys,
			"error": lastError,
			"speed": g.Speed,
		}
		select {
		case <-ctx.Done():
			return
		case keys, more = <-runes:
			if !more {
				return
			}
			if cmd, ok := KeyToCmd(keys...); ok {
				if err := g.RunCommand(cmd); err != nil {
					lastError = err
				}
			}
		case <-ticker.C:
			ticker.Reset(time.Duration(g.Speed))
			if err := g.Advance(); err != nil {
				// lastError = errors.Wrap(err, "GAME OVER!")
				cancel()
				showScreen(screens.GameOver)
			}
			// echo("advanced game", "step", game.Steps, "current", game.CurrentTile, "next", game.NextTile)
			if g.MaxSteps > 0 && g.Steps > g.MaxSteps {
				lastError = errors.New("GAME END!")
				cancel()
			}
		}
	}
}
