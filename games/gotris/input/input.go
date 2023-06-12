// Character input module defining a shared Key interface.
// Use this interface for handling key presses across packages.
// This package also implements an AwaitInput method.
package input

import (
	"context"
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

// safe test methods (works on nil)

func (k *Input) IsMovement() bool { return k != nil && k.flags.IsMovement() }
func (k *Input) IsAlt() bool      { return k != nil && k.flags.IsAlt() }
func (k *Input) IsText() bool     { return k != nil && k.rune != 0 }
func (k *Input) IsEmpty() bool    { return k == nil || k.key == KeyNone }

// unafe accessors to fields (will panic if Input is nil)

func (k *Input) Key() Key     { return k.key }
func (k *Input) Flags() Flag  { return k.flags }
func (k *Input) Rune() rune   { return k.rune }
func (k *Input) Text() string { return string(k.rune) }

// Await waits for user input or a given timeout.
func Await(ctx context.Context, input <-chan Input, timeout time.Duration) <-chan Input {
	ch := make(chan Input, 1)

	go func() {
		defer close(ch)

		// timeout zero means no time timeout, but wait forever (or parent context)
		if timeout > 0 {
			ctx2, cancel := context.WithTimeout(ctx, timeout)
			defer cancel()
			ctx = ctx2
		}

		// nil input means just wait for the time timeout or the parent context
		if input == nil {
			<-ctx.Done()
			return
		}

		// wait for input or timeout or parent context
		select {
		case <-ctx.Done():
			return
		case k := <-input:
			ch <- k
		}
	}()

	return ch
}
