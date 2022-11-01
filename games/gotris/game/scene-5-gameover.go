package game

import "github.com/ubntc/go/games/gotris/game/scenes"

func (g *Game) gameOver() {
	g.ShowScene(scenes.New(scenes.TitleGameOver), g.GameOverScreenDuration)
}
