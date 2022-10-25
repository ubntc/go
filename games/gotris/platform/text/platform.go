package text

import (
	"context"
	"os"
	"strings"

	"github.com/ubntc/go/games/gotris/game"
	"github.com/ubntc/go/games/gotris/game/scenes"
	"github.com/ubntc/go/games/gotris/rendering/text"
	"github.com/ubntc/go/games/gotris/rendering/text/modes"
	"github.com/ubntc/go/games/gotris/rendering/text/textscenes"
	"github.com/ubntc/go/games/gotris/terminal"
)

// Platform implements game rendering and input handling for the game,
// using the independent text rendering and terminal packages.
type Platform struct {
	terminal.Terminal
	modeMan *modes.ModeManager
}

func NewPlatform() *Platform {
	return &Platform{
		*terminal.New(os.Stdout),
		text.ModeMan(),
	}
}

func (p *Platform) Run(ctx context.Context) { <-ctx.Done() }

var DEBUG = os.Getenv("DEBUG") != ""

func (p *Platform) Render(g *game.Game) {
	p.Clear()
	p.Print(strings.Join(text.Render(g), "\r\n"))
}

func (p *Platform) RenderScene(scene *scenes.Scene) {
	var screen string
	switch scene.Name {
	case scenes.Welcome:
		screen = textscenes.WelcomeScreen.Menu(scene.Options, scene.Descriptions, scene.Options[scene.CurrentOption])
	case scenes.Options:
		screen = textscenes.OptionScreen.Menu(scene.Options, scene.Descriptions, scene.Options[scene.CurrentOption])
	case scenes.Controls:
		screen = textscenes.Controls
	case scenes.GameOver:
		screen = textscenes.GameOver
	}
	p.Clear()
	lines := strings.Split(string(screen), "\n")
	p.Print(strings.Join(lines, "\r\n    "))
}

func (p *Platform) RenderMessage(text string) {
	if DEBUG {
		lines := strings.Split("\n"+text, "\n")
		p.Print(strings.Join(lines, "\r\n"))
	}
}

func (p *Platform) SetRenderingMode(mode string) error {
	p.modeMan.SetModeByName(mode)
	return nil
}

func (p *Platform) RenderingModes() (names []string, currentMode int) {
	return p.modeMan.ModeNames(), p.modeMan.ModeIndex()
}

func (p *Platform) RenderingInfo(name string) string {
	return p.modeMan.ModeInfo(name)
}
