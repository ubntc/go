package game_test

import (
	"context"
	"testing"

	"github.com/ubntc/go/games/gotris/game"
	"github.com/ubntc/go/games/gotris/input"
)

type Platform struct{}

func (r *Platform) Render(game *game.Game)    {}
func (r *Platform) RenderScreen(text string)  {}
func (r *Platform) RenderMessage(text string) {}
func (r *Platform) CaptureInput(context.Context) (<-chan input.Key, func(), error) {
	return nil, nil, nil
}

func (p *Platform) SetRenderingMode(string) error   { return nil }
func (p *Platform) RenderingModes() ([]string, int) { return nil, 0 }
func (p *Platform) RenderingInfo(string) string     { return "" }

func TestLoop(t *testing.T) {
	game.NewGame(game.TestRules, &Platform{}).Run(game.CaptureOff)
}
