package text

import (
	"context"
	"os"
	"strings"

	"github.com/ubntc/go/games/gotris/common/labels"
	"github.com/ubntc/go/games/gotris/common/options"
	"github.com/ubntc/go/games/gotris/common/platform"
	"github.com/ubntc/go/games/gotris/terminal"
	"github.com/ubntc/go/games/gotris/ui/text/modes"
	txt "github.com/ubntc/go/games/gotris/ui/text/textscenes"
)

// TextUI implements game rendering and input handling for the game,
// using text rendering functions and the terminal package.
type TextUI struct {
	terminal.Terminal
	modeMan *modes.ModeManager
	options options.Options
}

func NewTextUI() *TextUI {
	mm := ModeMan()
	opts := renderingOptions{
		MemStore: *options.NewMemStore(
			mm.ModeNames(),
			mm.ModeDescs(),
		),
	}
	p := &TextUI{*terminal.New(os.Stdout), mm, &opts}
	opts.p = p
	opts.Set(mm.ModeIndex())
	return p
}

func (p *TextUI) Run(ctx context.Context) { <-ctx.Done() }

var DEBUG = os.Getenv("DEBUG") != ""

func (p *TextUI) Render(g platform.Game) {
	p.Clear()
	p.Print(strings.Join(Render(g), "\r\n"))
}

func (p *TextUI) RenderScene(scene platform.Scene) {
	var screen string
	opt := scene.Options()
	switch scene.Name() {
	case labels.TitleWelcome:
		screen = txt.WelcomeScreen.Menu(opt.List(), opt.Descs(), opt.GetName())
	case labels.TitleOptions:
		screen = txt.OptionScreen.Menu(opt.List(), opt.Descs(), opt.GetName())
	case labels.TitleControls:
		screen = txt.Controls
	case labels.TitleGameOver:
		screen = txt.GameOver
	}
	p.Clear()
	lines := strings.Split(string(screen), "\n")
	p.Print(strings.Join(lines, "\r\n    "))
}

func (p *TextUI) ShowMessage(text string) {
	if DEBUG {
		lines := strings.Split("\n"+text, "\n")
		p.Print(strings.Join(lines, "\r\n"))
	}
}

func (p *TextUI) SetRenderingMode(mode string) error {
	p.modeMan.SetModeByName(mode)
	return nil
}

func (p *TextUI) Options() options.Options {
	return p.options
}
