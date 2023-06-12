package game_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ubntc/go/games/gotris/game"
)

func TestGame_advance(t *testing.T) {
	assert := assert.New(t)
	g := game.NewGame(game.TestConfig, &Platform{})
	g.DisableInput()
	_ = g.Advance()
	assert.NotNil(g.CurrentTile)
	assert.NotNil(g.NextTile)
	for i := g.BoardSize.H; i > 0; i-- {
		_ = g.Advance()
	}
	assert.Greater(len(g.Board), 1)
}
