package game

import (
	"context"

	"github.com/ubntc/go/games/gotris/game/scenes"
)

func (g *Game) showWelcomeScreen(ctx context.Context) {
	idx := 0
	welcome := scenes.NewWelcomeMenu()

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		welcome.CurrentOption = idx
		key := g.ShowScene(welcome, 0)
		cmd, _ := KeyToMenuCmd(key)
		switch cmd {
		case CmdHelp:
			g.showScreen(scenes.Controls, 0)
		case CmdMenuDown, CmdMenuRight:
			idx = (idx + 1) % len(welcome.Options)
		case CmdMenuUp, CmdMenuLeft:
			idx = (idx + len(welcome.Options) - 1) % len(welcome.Options)
		case CmdMenuSelect:
			switch welcome.Options[idx] {
			case scenes.START:
				if err := g.GameLoop(ctx); err != nil {
					g.showScreen(scenes.GameOver, g.GameOverScreenDuration)
				}
			case scenes.OPTIONS:
				g.showOptions()
			case scenes.CONTROLS:
				g.showScreen(scenes.Controls, 0)
			case scenes.QUIT:
				return
			}
		case CmdEmpty, CmdQuit:
			return
		}
	}
}
