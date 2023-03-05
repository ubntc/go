package game

import (
	"context"
	"fmt"
	"time"

	"github.com/ubntc/go/games/gotris/common/labels"
	"github.com/ubntc/go/games/gotris/common/options"
	scenes "github.com/ubntc/go/games/gotris/common/platform"
	cmd "github.com/ubntc/go/games/gotris/game/controls"
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
	welcome := scenes.NewMenu(
		labels.TitleWelcome,
		options.NewMemStore([]string{labels.START, labels.OPTIONS, labels.CONTROLS, labels.QUIT}, nil),
	)
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

		if result := cmd.HandleOptionsCmd(c, opts); result == cmd.HandleResultSelectionFinished {
			switch opts.GetName() {
			case labels.START:
				g.Init()
				if err := g.GameLoop(ctx); err != nil {
					g.gameOver()
				}
			case labels.OPTIONS:
				g.showOptions()
			case labels.CONTROLS:
				g.showHelp()
			case labels.QUIT:
				return
			}
		}
	}
}
