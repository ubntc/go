package rendering

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ubntc/go/games/gotris/game"
)

func TestRenderer(t *testing.T) {
	assert := assert.New(t)

	g := game.NewGame(game.TestRules)
	step := 0
	for {
		step += 1
		out := Render(g)
		t.Log("text rendering, step", step, "\n"+strings.Join(out, "\n"))
		assert.Len(out, g.BoardSize.Height+3)
		if step == 10 {
			break
		}
		if err := g.Advance(); err != nil {
			t.Log("game over", err)
			break
		}
	}
}
