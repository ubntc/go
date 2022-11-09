// Character input module defining a shared Key interface.
// Use this interface for handling key presses across packages.
// This package also implements an AwaitInput method.
package input_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ubntc/go/games/gotris/input"
)

func TestIsMovement(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		name       string
		key        *input.Input
		wantIsMove bool
		wantKey    input.Key
	}{
		{"up", input.NewFromRune(65, input.FlagMove), true, input.KeyUp},
		{"b1", input.NewFromRune('Y', 0), false, input.KeyButton1},
		{"q", input.NewFromRune('q', 0), false, input.KeyQuit},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(tt.wantIsMove, tt.key.IsMovement(), "expect same movement flag")
			assert.Equal(tt.wantKey, tt.key.Key(), "expect same key")
		})
	}
}
