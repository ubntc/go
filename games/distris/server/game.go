package server

import (
	"context"
	"encoding/json"

	"github.com/ServiceWeaver/weaver"

	"github.com/ubntc/go/games/gotris/common/platform"
	gotris "github.com/ubntc/go/games/gotris/game"
	"github.com/ubntc/go/games/gotris/game/rules"
	ui "github.com/ubntc/go/games/gotris/ui/text"
)

// Game component.
type Game interface {
	TogglePause(_ context.Context) error
	Play(_ context.Context) error
	Get(_ context.Context) (*GameState, error)
}

// Implementation of the Game component.
type game struct {
	weaver.Implements[Game]

	gg gotris.Game
}

type GameState struct {
	Game *platform.Game
}

func (s *GameState) MarshalBinary() (data []byte, err error) {
	return json.Marshal(s.Game)
}

func (s *GameState) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &s.Game)
}

func (g *game) Init(ctx context.Context) error {
	root := weaver.Init(ctx)

	g.gg = *gotris.NewGame(rules.DefaultRules, ui.NewTextUI())

	g.Logger().Info("init game", "root", root)
	return nil
}

func (g *game) TogglePause(_ context.Context) error {
	return nil
}

func (g *game) Play(_ context.Context) error {
	return g.gg.Advance()
}

func (g *game) Get(_ context.Context) (*GameState, error) {
	return &GameState{&g.gg.Game}, nil
}
