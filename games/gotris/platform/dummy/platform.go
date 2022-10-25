package fyne

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/ubntc/go/games/gotris/game"
)

var DEBUG = os.Getenv("DEBUG") != ""

// Platform implements game rendering and input handling for the game.
type Platform struct{}

func NewPlatform() *Platform {
	return &Platform{}
}

func (p *Platform) Run(ctx context.Context) {
	<-ctx.Done()
}

func echo(s ...interface{}) {
	text := fmt.Sprintln(s...)
	lines := strings.Split("\n"+text, "\n")
	fmt.Println(strings.Join(lines, "\r\n"))
}

func (p *Platform) Render(g *game.Game) {
	echo("Board", g.Board)
	echo("Score", g.Score)
}

func (p *Platform) RenderScene(text string)   { echo(text) }
func (p *Platform) RenderMessage(text string) { echo(text) }

func (p *Platform) SetRenderingMode(mode string) error {
	return nil
}

func (p *Platform) RenderingModes() (names []string, currentMode int) {
	return []string{"default"}, 0
}

func (p *Platform) RenderingInfo(name string) string {
	return "default mode"
}
