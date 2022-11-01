package game

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
	cmd "github.com/ubntc/go/games/gotris/game/controls"
)

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
			g.platform.ShowMessage(msg)
		}

		select {
		case <-ctx.Done():
			return nil
		case key, more := <-g.input:
			if !more {
				return nil
			}
			msg = fmt.Sprintf("key(%v, %v, %v)", key.Mod(), key.Rune(), key.Runes())
			if c, arg := cmd.KeyToCmd(key); c != cmd.Empty {
				msg += fmt.Sprintf(", cmd: %s, arg:%s", c, arg)
				lastErr = g.RunCommand(c, arg)
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