package main

import (
	"os"
	"strings"

	"github.com/ubntc/go/games/gotris/game"
	"github.com/ubntc/go/games/gotris/rendering"
	"github.com/ubntc/go/games/gotris/rendering/modes"
	"github.com/ubntc/go/games/gotris/terminal"
)

// Platform implements game rendering and input handling for the game,
// using the independent text rendering and terminal packages.
type Platform struct {
	terminal.Terminal
	modeMan *modes.ModeManager
}

var DEBUG = os.Getenv("DEBUG") != ""

func (p *Platform) Render(g *game.Game) {
	p.Clear()
	p.Print(strings.Join(rendering.Render(g), "\r\n"))
}

func (p *Platform) RenderScreen(text string) {
	p.Clear()
	lines := strings.Split(text, "\n")
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

func main() {
	p := Platform{
		*terminal.New(os.Stdout),
		rendering.ModeMan(),
	}
	game.NewGame(game.DefaultRules, &p).Run(game.CaptureOn)
}
