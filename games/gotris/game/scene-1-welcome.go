package game

import (
	"context"
	"fmt"

	cmd "github.com/ubntc/go/games/gotris/game/controls"
	"github.com/ubntc/go/games/gotris/game/scenes"
)

func (g *Game) handleHelpAndQuit(ctx context.Context, c cmd.Cmd) (quit, ok bool) {
	switch c {
	case cmd.Help:
		g.showHelp(ctx)
		return false, true // do not quit + handled command
	case cmd.Empty, cmd.Quit:
		return true, true // do quit + handled command
	default:
		return false, false // do not quit + did not handle command
	}
}

func (g *Game) showWelcome(ctx context.Context) {
	welcome := scenes.NewWelcomeMenu()
	opts := welcome.Options()

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		key := g.showScene(ctx, welcome, 0)
		fmt.Println(key)
		c, _ := cmd.KeyToMenuCmd(key)

		if quit, ok := g.handleHelpAndQuit(ctx, c); ok {
			if quit {
				return
			}
			continue
		}

		if result := cmd.HandleOptionsCmd(c, opts); result == cmd.HandleResultSelectionFinished {
			switch opts.GetName() {
			case scenes.START:
				if err := g.GameLoop(ctx); err != nil {
					g.gameOver(ctx)
				}
			case scenes.OPTIONS:
				g.showOptions(ctx)
			case scenes.CONTROLS:
				g.showHelp(ctx)
			case scenes.QUIT:
				return
			}
		}
	}
}
