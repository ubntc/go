// Character input module defining a shared Key interface.
// Use this interface for handling key presses across packages.
// This package also implements an AwaitInput method.
package input

import (
	"time"
)

type Input struct {
	key   Key
	flags Flag
	rune  rune
}

var Empty = New(KeyNone, FlagNone)

func New(k Key, flags Flag) *Input {
	return &Input{k, flags, 0}
}

func (k *Input) IsMovement() bool { return k.flags.IsMovement() }
func (k *Input) IsAlt() bool      { return k.flags.IsAlt() }
func (k *Input) IsText() bool     { return k.rune != 0 }
func (k *Input) Key() Key         { return k.key }
func (k *Input) Flags() Flag      { return k.flags }
func (k *Input) Rune() rune       { return k.rune }
func (k *Input) Text() string     { return string(k.rune) }

// AwaitInput waits for user input or a given timeout.
func AwaitInput(input <-chan *Input, timeout time.Duration) *Input {
	switch {
	case input == nil:
		time.Sleep(timeout)
	case timeout == 0:
		return <-input
	case timeout != 0:
		select {
		case k := <-input:
			return k
		case <-time.After(timeout):
		}
	}
	return Empty
}
