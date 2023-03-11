package server

import (
	"context"

	"github.com/ServiceWeaver/weaver"
	"github.com/ubntc/go/games/distris/api/command"
)

type CommandStatus string

const (
	CommandInvalid    = CommandStatus("command is invalid")
	CommandSuccessful = CommandStatus("command was processed successfully")
	CommandFailed     = CommandStatus("command was processed with errors")
)

// Router component.
type Router interface {
	Send(context.Context, command.Command) (CommandStatus, error)
}

// Implementation of the router component.
type router struct {
	weaver.Implements[Router]

	game Game
}

func (r *router) Init(ctx context.Context) error {
	root := weaver.Init(ctx)

	// Get games to be played.
	game, err := weaver.Get[Game](root)
	if err != nil {
		return err
	}
	r.game = game
	return nil
}

func (r *router) Send(ctx context.Context, cmd command.Command) (CommandStatus, error) {
	var err error
	var game *GameState

	switch cmd {
	case command.EMPTY:
		return CommandSuccessful, nil
	case command.PLAY:
		err = r.game.Play(ctx)
	case command.PAUSE:
		err = r.game.TogglePause(ctx)
	case command.GET:
		game, err = r.game.Get(ctx)
	default:
		return CommandInvalid, nil
	}

	if err != nil {
		return CommandFailed, err
	}

	if game != nil {
		str, err := game.MarshalBinary()
		return CommandStatus(string(str)), err
	}

	return CommandSuccessful, nil
}
