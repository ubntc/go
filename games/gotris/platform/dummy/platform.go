package dummy

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/ubntc/go/games/gotris/game"
	"github.com/ubntc/go/games/gotris/game/options"
	"github.com/ubntc/go/games/gotris/game/scenes"
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
	echo("Tile", g.CurrentTile)
}

func (p *Platform) RenderScene(scene *scenes.Scene) { echo(scene.Name()) }
func (p *Platform) ShowMessage(text string)         { echo(text) }
func (p *Platform) Options() options.Options        { return options.NewMemStore(nil, nil) }
