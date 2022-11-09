package main

import (
	"context"
	"flag"

	"github.com/ubntc/go/games/gotris/game"
	"github.com/ubntc/go/games/gotris/game/rules"
	"github.com/ubntc/go/games/gotris/platform/dummy"
	"github.com/ubntc/go/games/gotris/platform/fyne"
	"github.com/ubntc/go/games/gotris/platform/text"
)

const (
	PlatformText  = "text"
	PlatformFyne  = "fyne"
	PlatformDummy = "dummy"
)

var platform = flag.String("platform", "text", "name of the game platform to use")

func main() {
	flag.Parse()
	var p game.Platform
	switch *platform {
	case PlatformFyne:
		p = fyne.NewPlatform()
	case PlatformText:
		p = text.NewPlatform()
	case PlatformDummy:
		p = dummy.NewPlatform()
	}
	g := game.NewGame(rules.DefaultRules, p)

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
