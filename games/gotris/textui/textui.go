package textui

import (
	"context"
	"os"
	"strings"

	"github.com/ubntc/go/games/gotris/common/options"
	"github.com/ubntc/go/games/gotris/common/platform"
	"github.com/ubntc/go/games/gotris/terminal"
	"github.com/ubntc/go/games/gotris/textui/modes"
)

var DEBUG = os.Getenv("DEBUG") != ""

// TextUI implements game rendering and input handling for the game,
// using text rendering functions and the terminal package.
type TextUI struct {
	terminal.Terminal
	modeMan *modes.ModeManager
	options options.Options
}

func NewTextUI() *TextUI {
	mm := modes.NewModeManager()
	opts := options.NewMemStore(
		mm.ModeNames(),
		mm.ModeDescs(),
	)
	opts.Set(mm.ModeIndex())

	return &TextUI{
		Terminal: *terminal.NewTerminal(os.Stdout),
		modeMan:  mm,
		options:  opts,
	}
}

func (ui *TextUI) Run(ctx context.Context) {
	for {
		select {
		case <-ui.options.Changed():
			ui.modeMan.SetModeByName(ui.options.GetName())
		case <-ctx.Done():
			return
		}
	}
}

func (ui *TextUI) Render(g platform.Game) {
	ui.Clear()
	ui.Print(strings.Join(ui.RenderGame(g), "\r\n"))
}

func (ui *TextUI) ShowMessage(text string) {
	if DEBUG {
		lines := strings.Split("\n"+text, "\n")
		ui.Print(strings.Join(lines, "\r\n"))
	}
}

func (ui *TextUI) Options() options.Options {
	return ui.options
}
