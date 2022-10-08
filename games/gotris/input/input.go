// Character input module defining a shared Key interface.
// Use this interface for handling key presses across packages.
// This package also implements an AwaitInput method.
package input

import (
	"time"
)

type Mod int

const (
	ModNone  Mod = 0
	ModShift Mod = 1
	ModAlt   Mod = 2
	ModCtrl  Mod = 4
	ModMove  Mod = 8
)

type Key interface {
	Rune() rune
	Mod() Mod
	Runes() []rune
}

// AwaitInput waits for user input or a given timeout.
func AwaitInput(input <-chan Key, timeout time.Duration) (key Key) {
	switch {
	case input == nil:
		time.Sleep(timeout)
	case timeout == 0:
		key = <-input
	case timeout != 0:
		select {
		case key = <-input:
		case <-time.After(timeout):
		}
	}
	return
}
