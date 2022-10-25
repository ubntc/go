package main

import (
	"context"
	"flag"

	"github.com/ubntc/go/games/gotris/game"
	"github.com/ubntc/go/games/gotris/platform/fyne"
	"github.com/ubntc/go/games/gotris/platform/text"
)

var platform = flag.String("platform", "fyne", "name of the game platform to use")

func main() {
	flag.Parse()
	var p game.Platform
	switch *platform {
	case "fyne":
		p = fyne.NewPlatform()
	case "text":
		p = text.NewPlatform()
	}
	g := game.NewGame(game.DefaultRules, p)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// start the game loop in the background
	go func() {
		g.Run(ctx)
		cancel()
	}()

	// handover main thead to the platform
	p.Run(ctx)
}
