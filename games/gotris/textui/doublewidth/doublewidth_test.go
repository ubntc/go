package doublewidth_test

import (
	"testing"

	"github.com/ubntc/go/games/gotris/textui/doublewidth"
)

func TestTextToBlock(t *testing.T) {
	dw := doublewidth.New()
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
			if got := dw.TextToBlock(tt.text); got != tt.want {
				t.Errorf("TextToBlock() = %v, want %v", got, tt.want)
			}
		})
	}
}
