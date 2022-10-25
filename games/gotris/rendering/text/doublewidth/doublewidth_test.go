package doublewidth_test

import (
	"testing"

	"github.com/ubntc/go/games/gotris/rendering/text/doublewidth"
)

func TestTextToBlock(t *testing.T) {
	art := doublewidth.New()
	puncts := ",.-;:_#'+*^!\"$%&/()=?"
	tests := []struct {
		name string
		text string
		want string
	}{
		{"numbers", "0123456789", "0123456789"},
		{"abc", "abcABCxyzXYZ", "abcABCxyzXYZ"},
		{"puncts", puncts, puncts},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := art.TextToBlock(tt.text); got != tt.want {
				t.Errorf("TextToBlock() = %v, want %v", got, tt.want)
			}
		})
	}
}
