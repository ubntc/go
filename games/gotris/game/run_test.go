package game_test

import (
	"context"
	"testing"

	"github.com/ubntc/go/games/gotris/game"
)

type Renderer struct{}

func (r *Renderer) Render(game *game.Game) {}
func (r *Renderer) RenderText(text string) {}
func (r *Renderer) CaptureInput(context.Context) (<-chan []rune, func(), error) {
	return nil, nil, nil
}

func TestLoop(t *testing.T) {
	game.Run(game.TestRules, &Renderer{}, game.CaptureOff)
}
