package main

import (
	"os"
	"strings"

	"github.com/ubntc/go/games/gotris/game"
	"github.com/ubntc/go/games/gotris/rendering"
	"github.com/ubntc/go/games/gotris/terminal"
)

// Platform implements game rendering and input handling for the game,
// using the independent text rendering and terminal packages.
type Platform struct {
	terminal.Terminal
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

func main() {
	p := Platform{*terminal.New(os.Stdout)}
	game.NewGame(game.DefaultRules, &p).Run(game.CaptureOn)
}
