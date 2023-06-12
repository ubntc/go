package game

import (
	"context"

	"github.com/ubntc/go/games/gotris/game/scenes"
)

func (g *Game) gameOver(ctx context.Context) {
	g.showScene(ctx, scenes.New(scenes.TitleGameOver), g.GameOverScreenDuration)
}
