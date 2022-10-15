// Character input module defining a shared Key interface.
// Use this interface for handling key presses across packages.
// This package also implements an AwaitInput method.
package input

import (
	"testing"
)

type K struct {
	mod Mod
}

func (k *K) Mod() Mod      { return k.mod }
func (k *K) Rune() rune    { return 0 }
func (k *K) Runes() []rune { return nil }

func TestIsMovement(t *testing.T) {
	tests := []struct {
		name string
		key  K
		want bool
	}{
		{"space", K{mod: ModNone}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsMovement(&tt.key); got != tt.want {
				t.Errorf("IsMovement() = %v, want %v", got, tt.want)
			}
		})
	}
}
