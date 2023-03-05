package game

import (
	"github.com/ubntc/go/games/gotris/common/labels"
	"github.com/ubntc/go/games/gotris/common/platform"
)

func (g *Game) showHelp() {
	g.ShowScene(platform.NewScene(labels.TitleControls), 0)
}
