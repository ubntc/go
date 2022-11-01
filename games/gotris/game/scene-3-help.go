package game

import "github.com/ubntc/go/games/gotris/game/scenes"

func (g *Game) showHelp() {
	g.ShowScene(scenes.New(scenes.TitleControls), 0)
}
