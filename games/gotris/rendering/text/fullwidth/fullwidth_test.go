package fullwidth_test

import (
	"testing"

	"github.com/ubntc/go/games/gotris/rendering/text/fullwidth"
)

func TestFullwidthTextTranslation(t *testing.T) {
	art := fullwidth.New()
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
			if got := art.TextToBlock(tt.text); got != tt.want {
				t.Errorf("TextToBlock() = %v, want %v", got, tt.want)
			}
		})
	}
}
