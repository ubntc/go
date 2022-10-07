package rendering

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ubntc/go/games/gotris/game"
)

func TestRenderer(t *testing.T) {
	assert := assert.New(t)

	g := game.NewGame(game.TestRules, nil)
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

func TestTextToBlock(t *testing.T) {
	tests := []struct {
		name string
		text string
		want string
	}{
		{"numbers", "0123456789", "０１２３４５６７８９"},
		{"punctuation", ".,:;!?", "．，：；！？"},
		{"math", "+-*=/", "＋－*＝／"},
		{"common signs", "$%&@#'^~_", "＄％＆＠＃＇＾～＿"},
		{"braces", "()[]{}|", "（）［］｛｝｜"},
		{"space", " ", "　"},
		{"empty", "", ""},
		{"quote", "\"", "＂"},
		{"5 spaces", "     ", "　　　　　"},
		{"Gotris", "Gotris", "Ｇｏｔｒｉｓ"},
		{"Game Over", "Game Over", "Ｇａｍｅ　Ｏｖｅｒ"},
		{"GOTRIS", "GOTRIS", "ＧＯＴＲＩＳ"},
		{"GAME OVER", "GAME OVER", "ＧＡＭＥ　ＯＶＥＲ"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TextToBlock(tt.text); got != tt.want {
				t.Errorf("TextToBlock() = %v, want %v", got, tt.want)
			}
		})
	}
}
