package text

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ubntc/go/games/gotris/game"
	"github.com/ubntc/go/games/gotris/game/rules"
)

func TestRendering(t *testing.T) {
	assert := assert.New(t)

	g := game.NewGame(rules.TestRules, nil)
	g.CaptureInput = false
	step := 0
	for {
		step += 1
		out := Render(g.Game)
		t.Log("text rendering, step", step, "\n"+strings.Join(out, "\n"))
		assert.Len(out, g.BoardSize.H+3)
		if step == 10 {
			break
		}
		if err := g.Advance(); err != nil {
			t.Log("game over", err)
			break
		}
	}
}
