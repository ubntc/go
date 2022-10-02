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

func (p *Platform) Render(g *game.Game) {
	p.Clear()
	p.Print(strings.Join(rendering.Render(g), "\r\n"))
}

func main() {
	p := Platform{*terminal.New(os.Stdout)}
	p.Print("starting Gotris")
	game.Run(game.DefaultRules, &p, game.CaptureOn)
	p.Println("\nGotris stopped")
}
