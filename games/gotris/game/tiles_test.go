package game_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ubntc/go/games/gotris/game"
)

func TestPoint(t *testing.T) {
	assert := assert.New(t)

	a := game.Point{X: 1, Y: 1}
	b := game.Point{X: 1, Y: 1}

	assert.True(a == b, "ensure points can be compared")
}
