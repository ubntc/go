package game_test

import (
	"context"
	"testing"
	"time"

	"github.com/ubntc/go/games/gotris/common/input"
	"github.com/ubntc/go/games/gotris/common/options"
	"github.com/ubntc/go/games/gotris/common/platform"
	"github.com/ubntc/go/games/gotris/game"
	"github.com/ubntc/go/games/gotris/game/rules"
)

type Platform struct{}

func (r *Platform) Render(game platform.Game)        {}
func (r *Platform) RenderScene(scene platform.Scene) {}
func (r *Platform) ShowMessage(text string)          {}
func (r *Platform) CaptureInput(context.Context) (<-chan *input.Input, func(), error) {
	return nil, nil, nil
}
func (p *Platform) Options() options.Options { return nil }

func (p *Platform) SetRenderingMode(string) error   { return nil }
func (p *Platform) RenderingModes() ([]string, int) { return nil, 0 }
func (p *Platform) RenderingInfo(string) string     { return "" }
func (p *Platform) Run(ctx context.Context)         { <-ctx.Done() }

func TestLoop(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	p := &Platform{}
	g := game.NewGame(rules.TestRules, p)
	g.CaptureInput = false
	go g.Run(ctx)
	p.Run(ctx)
}
