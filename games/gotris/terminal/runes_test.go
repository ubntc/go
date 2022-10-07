package terminal

import "testing"

func Test_handleRune(t *testing.T) {
	tests := []struct {
		name          string
		historyLength int
		newRune       rune
		want          action
	}{
		{"quit", 0, 'q', actionQuit},
		{"end of text ctrl-c 3", 0, 3, actionQuit},
		{"end of xmit ctrl-d 4", 0, 4, actionQuit},

		{"escape code start 27", 0, 27, actionAppendEscape},
		{"escape code navigation keys 91", 1, 91, actionAppendMovement},
		{"escape final code arrow up", 2, 65, actionAppendAndSend},
		{"escape final code arrow down", 2, 66, actionAppendAndSend},
		{"escape final code arrow right", 2, 67, actionAppendAndSend},
		{"escape final code arrow left", 2, 68, actionAppendAndSend},

		{"char 'x'", 0, 'x', actionAppendAndSend},
		{"char 'x' after 1 ctrl", 1, 'x', actionAppendAndSend},
		{"char 'x' after 2 ctrl", 2, 'x', actionAppendAndSend},

		{"buffer too long", 10, 'x', actionAppendAndSend},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := handleRune(tt.historyLength, tt.newRune); got != tt.want {
				t.Errorf("handleRune() = %v, want %v", got, tt.want)
			}
		})
	}
}
