package text

import (
	"context"
	"os"
	"strings"

	"github.com/ubntc/go/games/gotris/game"
	"github.com/ubntc/go/games/gotris/game/scenes"
	"github.com/ubntc/go/games/gotris/rendering/text"
	"github.com/ubntc/go/games/gotris/rendering/text/modes"
	txt "github.com/ubntc/go/games/gotris/rendering/text/textscenes"
	"github.com/ubntc/go/games/gotris/terminal"
)

// Platform implements game rendering and input handling for the game,
// using the independent text rendering and terminal packages.
type Platform struct {
	terminal.Terminal
	modeMan *modes.ModeManager
	options scenes.Options
}

func NewPlatform() *Platform {
	mm := text.ModeMan()
	opts := RenderingOptions{
		SceneOptions: scenes.SceneOptions{
			Options:      mm.ModeNames(),
			Descriptions: mm.ModeDescs(),
		},
	}
	p := &Platform{*terminal.New(os.Stdout), mm, &opts}
	opts.p = p
	opts.Set(mm.ModeIndex())
	return p
}

func (p *Platform) Run(ctx context.Context) { <-ctx.Done() }

var DEBUG = os.Getenv("DEBUG") != ""

func (p *Platform) Render(g *game.Game) {
	p.Clear()
	p.Print(strings.Join(text.Render(g), "\r\n"))
}

func (p *Platform) RenderScene(scene scenes.Scene) {
	var screen string
	opt := scene.Options()
	switch scene.Name() {
	case scenes.TitleWelcome:
		screen = txt.WelcomeScreen.Menu(opt.List(), opt.Descs(), opt.GetName())
	case scenes.TitleOptions:
		screen = txt.OptionScreen.Menu(opt.List(), opt.Descs(), opt.GetName())
	case scenes.TitleControls:
		screen = txt.Controls
	case scenes.TitleGameOver:
		screen = txt.GameOver
	}
	p.Clear()
	lines := strings.Split(string(screen), "\n")
	p.Print(strings.Join(lines, "\r\n    "))
}

func (p *Platform) ShowMessage(text string) {
	if DEBUG {
		lines := strings.Split("\n"+text, "\n")
		p.Print(strings.Join(lines, "\r\n"))
	}
}

func (p *Platform) SetRenderingMode(mode string) error {
	p.modeMan.SetModeByName(mode)
	return nil
}

func (p *Platform) Options() scenes.Options {
	return p.options
}
