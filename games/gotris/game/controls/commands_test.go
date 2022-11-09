package controls_test

import (
	"testing"

	cmd "github.com/ubntc/go/games/gotris/game/controls"
	"github.com/ubntc/go/games/gotris/input"
)

func TestKeyToCmd(t *testing.T) {
	tests := []struct {
		name    string
		in      *input.Input
		wantCmd cmd.Cmd
		wantArg string
	}{
		{"drop", input.New(input.KeyButton3, 0), cmd.Drop, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCmd, gotOk := cmd.InputToCmd(tt.in)
			if gotCmd != tt.wantCmd {
				t.Errorf("KeyToCmd() gotCmd = %v, want %v", gotCmd, tt.wantCmd)
			}
			if gotOk != tt.wantArg {
				t.Errorf("KeyToCmd() gotOk = %v, want %v", gotOk, tt.wantArg)
			}
		})
	}
}
