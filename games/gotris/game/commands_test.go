package game_test

import (
	"testing"

	cmd "github.com/ubntc/go/games/gotris/game/controls"
	"github.com/ubntc/go/games/gotris/input"
)

type K struct {
	mod  input.Mod
	rune rune
}

func (k *K) Mod() input.Mod { return k.mod }
func (k *K) Rune() rune     { return k.rune }
func (k *K) Runes() []rune  { return []rune{k.rune} }

func TestKeyToCmd(t *testing.T) {
	tests := []struct {
		name    string
		key     K
		wantCmd cmd.Cmd
		wantArg string
	}{
		{"drop", K{0, ' '}, cmd.Drop, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCmd, gotOk := cmd.KeyToCmd(&tt.key)
			if gotCmd != tt.wantCmd {
				t.Errorf("KeyToCmd() gotCmd = %v, want %v", gotCmd, tt.wantCmd)
			}
			if gotOk != tt.wantArg {
				t.Errorf("KeyToCmd() gotOk = %v, want %v", gotOk, tt.wantArg)
			}
		})
	}
}
