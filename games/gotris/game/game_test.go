package game_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ubntc/go/games/gotris/game"
	"github.com/ubntc/go/games/gotris/game/rules"
)

func TestGame_advance(t *testing.T) {
	assert := assert.New(t)
	g := game.NewGame(rules.TestRules, &Platform{})
	g.CaptureInput = false
	g.Advance()
	assert.NotNil(g.CurrentTile)
	assert.NotNil(g.NextTile)
	for i := g.BoardSize.H; i > 0; i-- {
		g.Advance()
	}
	assert.Greater(len(g.Board), 1)
}
