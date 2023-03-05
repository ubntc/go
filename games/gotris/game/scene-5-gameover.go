package game

import (
	"github.com/ubntc/go/games/gotris/common/labels"
	"github.com/ubntc/go/games/gotris/common/platform"
)

func (g *Game) gameOver() {
	g.ShowScene(platform.NewScene(labels.TitleGameOver), g.GameOverScreenDuration)
}
