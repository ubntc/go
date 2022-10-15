package game

import (
	"context"
	"fmt"
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

	// do not capture input during tests
	if capture {
		ch, restore, err := g.platform.CaptureInput(ctx)
		if err != nil {
			return err
		}
		defer restore()
		g.input = ch
		g.GameOverScreenDuration = time.Second * 5
	}

	g.MainLoop(ctx)

	return nil
}

func (g *Game) MainLoop(ctx context.Context) {
	// redraw the game to show the last Game state after any previous screen
	defer g.platform.Render(g)

	idx := 0
	options := []string{screens.START, screens.OPTIONS, screens.CONTROLS, screens.QUIT}

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		key := g.ShowScreen(screens.WelcomeScreen.Menu(options, nil, options[idx]), 0)
		cmd, _ := KeyToMenuCmd(key)
		switch cmd {
		case CmdHelp:
			g.ShowScreen(screens.Controls, 0)
		case CmdMenuDown, CmdMenuRight:
			idx = (idx + 1) % len(options)
		case CmdMenuUp, CmdMenuLeft:
			idx = (idx + len(options) - 1) % len(options)
		case CmdMenuSelect:
			switch options[idx] {
			case screens.START:
				if err := g.GameLoop(ctx); err != nil {
					g.ShowScreen(screens.GameOver, g.GameOverScreenDuration)
				}
			case screens.OPTIONS:
				g.showOptions()
			case screens.CONTROLS:
				g.ShowScreen(screens.Controls, 0)
			case screens.QUIT:
				return
			}
		case CmdEmpty, CmdQuit:
			return
		}
	}
}

func (g *Game) GameLoop(ctx context.Context) error {
	var lastErr error
	var msg string
	ticker := time.NewTicker(g.TickTime)
	defer ticker.Stop()
	for {
		g.platform.Render(g)
		if lastErr != nil {
			msg = lastErr.Error()
		}
		if msg != "" {
			g.platform.RenderMessage(msg)
		}

		select {
		case <-ctx.Done():
			return nil
		case key, more := <-g.input:
			if !more {
				return nil
			}
			msg = fmt.Sprintf("key(%v, %v, %v)", key.Mod(), key.Rune(), key.Runes())
			if cmd, arg := KeyToCmd(key); cmd != CmdEmpty {
				msg += fmt.Sprintf(", cmd: %s, arg:%s", cmd, arg)
				lastErr = g.RunCommand(cmd, arg)
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
