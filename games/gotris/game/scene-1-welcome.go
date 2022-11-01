package game

import (
	"context"
	"fmt"
	"time"

	cmd "github.com/ubntc/go/games/gotris/game/controls"
	"github.com/ubntc/go/games/gotris/game/scenes"
)

func (g *Game) handleCommonCommand(c cmd.Cmd) (quit, ok bool) {
	switch c {
	case cmd.Help:
		g.showHelp()
	case cmd.Empty, cmd.Quit:
		quit = true
	default:
		return false, false
	}
	return quit, true
}

func hint[K any](v ...K) {
	fmt.Printf("%v\n", v)
	time.Sleep(time.Second)
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

		key := g.ShowScene(welcome, 0)
		c, _ := cmd.KeyToMenuCmd(key)

		if quit, ok := g.handleCommonCommand(c); ok {
			if quit {
				return
			}
			continue
		}

		if i, done, ok := cmd.HandleOptionsCmd(c, len(opts.List()), opts.Get()); ok {
			opts.Set(i)
			if !done {
				continue
			}
			switch opts.GetName() {
			case scenes.START:
				if err := g.GameLoop(ctx); err != nil {
					g.gameOver()
				}
			case scenes.OPTIONS:
				g.showOptions()
			case scenes.CONTROLS:
				g.showHelp()
			case scenes.QUIT:
				return
			}
		}
	}
}
