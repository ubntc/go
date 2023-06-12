package game

import (
	"context"

	"github.com/ubntc/go/games/gotris/game/scenes"
)

func (g *Game) showHelp(ctx context.Context) {
	g.showScene(ctx, scenes.New(scenes.TitleControls), 0)
}
