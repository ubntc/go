package fyne

import (
	"context"
	"image/color"
	"math/rand"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"github.com/ubntc/go/games/gotris/game"
	"github.com/ubntc/go/games/gotris/game/scenes"
)

var DEBUG = os.Getenv("DEBUG") != ""

// Platform implements game rendering and input handling for the game.
type Platform struct {
	app fyne.App
	wnd fyne.Window
	pix fyne.Canvas
}

func NewPlatform() *Platform {
	p := &Platform{
		app: app.New(),
	}
	p.wnd = p.app.NewWindow("Gotris")
	p.pix = p.wnd.Canvas()
	return p
}

func (p *Platform) Run(ctx context.Context) {
	go func() {
		<-ctx.Done()
		p.app.Quit()
	}()
	p.wnd.Resize(fyne.NewSize(400, 800))
	p.wnd.ShowAndRun()
}

func (p *Platform) ShowMessage(text string) {
	echo(text)
}

func (p *Platform) Options() scenes.Options {
	return nil
}

func (p *Platform) Render(g *game.Game) {
	R, G, B := rand.Int()%255, rand.Int()%255, rand.Int()%255
	blue := color.NRGBA{R: uint8(R), G: uint8(G), B: uint8(B), A: 255}
	rect := canvas.NewRectangle(blue)
	p.pix.SetContent(rect)
	echo("Board", g.Board)
	echo("Score", g.Score)
}

func (p *Platform) RenderScene(scene scenes.Scene) {
	content := canvas.NewText(scene.Name(), WHITE)
	p.pix.SetContent(content)
	echo(scene)
}
