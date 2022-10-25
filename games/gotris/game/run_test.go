package game_test

import (
	"context"
	"testing"
	"time"

	"github.com/ubntc/go/games/gotris/game"
	"github.com/ubntc/go/games/gotris/game/scenes"
	"github.com/ubntc/go/games/gotris/input"
)

type Platform struct{}

func (r *Platform) Render(game *game.Game)          {}
func (r *Platform) RenderScene(scene *scenes.Scene) {}
func (r *Platform) RenderMessage(text string)       {}
func (r *Platform) CaptureInput(context.Context) (<-chan input.Key, func(), error) {
	return nil, nil, nil
}

func (p *Platform) SetRenderingMode(string) error   { return nil }
func (p *Platform) RenderingModes() ([]string, int) { return nil, 0 }
func (p *Platform) RenderingInfo(string) string     { return "" }
func (p *Platform) Run(ctx context.Context)         { <-ctx.Done() }

func TestLoop(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	p := &Platform{}
	g := game.NewGame(game.TestRules, p)
	g.CaptureInput = false
	go g.Run(ctx)
	p.Run(ctx)
}
