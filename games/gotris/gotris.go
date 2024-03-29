package main

import (
	"context"

	"github.com/ubntc/go/games/gotris/game"
	"github.com/ubntc/go/games/gotris/game/rules"
	"github.com/ubntc/go/games/gotris/textui"
)

func main() {
	ui := textui.NewTextUI()
	g := game.NewGame(rules.DefaultRules, ui)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// start the game loop in the background
	go func() {
		g.Run(ctx)
		cancel()
	}()

	// handover main thead to the UI
	ui.Run(ctx)
}
