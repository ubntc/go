package terminal

import (
	"bufio"
	"context"
	"log"
	"os"

	"github.com/pkg/errors"
	"github.com/ubntc/go/games/gotris/input"
	xterm "golang.org/x/term"
)

var debug = os.Getenv("DEBUG") != ""

type key struct {
	rune  rune
	mod   input.Mod
	runes []rune
}

func (k *key) Rune() rune {
	return k.rune
}

func (k *key) Mod() input.Mod {
	return k.mod
}

func (k *key) Runes() []rune {
	return k.runes
}

func (t *Terminal) CaptureInput(ctx context.Context) (<-chan input.Key, func(), error) {
	keys := make(chan input.Key)
	stdin := int(t.stdin.Fd())

	state, err := xterm.MakeRaw(stdin)
	if err != nil {
		return nil, nil, errors.Wrap(err, "term.MakeRaw")
	}
	if state == nil {
		return nil, nil, errors.Wrap(errors.New("restore state must not be nil"), "term.MakeRaw")
	}
	t.HideCursor()

	restore := func() {
		t.ShowCursor()
		if err := xterm.Restore(stdin, state); err != nil {
			log.Fatalln(errors.Wrap(err, "term.Restore"))
		}
	}

	sendKey := func(r rune, mod input.Mod, runes []rune) {
		select {
		case keys <- &key{r, input.Mod(mod), runes}:
		default:
			// ignore new input if prev. input is not processed
		}
	}

	stopRequested := func() bool {
		select {
		case <-ctx.Done():
			return true
		default:
		}
		return false
	}

	go func() {
		defer func() {
			restore()
			close(keys)
		}()
		in := bufio.NewReader(os.Stdin)
		var buf []rune
		var mod input.Mod

		for {
			if stopRequested() {
				return
			}

			r, _, err := in.ReadRune()
			if err != nil {
				if debug {
					log.Println(err)
				}
				return
			}

			switch handleRune(len(buf), r) {
			case actionSendWithAltAsMovement:
				mod |= input.ModAlt | input.ModMove
				fallthrough
			case actionAppendAndSend:
				sendKey(r, mod, append(buf, r))
				mod = 0
				buf = nil
			case actionAppendAlt:
				mod |= input.ModAlt
			case actionAppendCtrl:
				mod |= input.ModCtrl
			case actionAppendShift:
				mod |= input.ModShift
			case actionAppendAltShift:
				mod |= input.ModAlt | input.ModShift
			case actionAppendCtrlShift:
				mod |= input.ModCtrl | input.ModShift
			case actionAppendMovement:
				mod |= input.ModMove
				fallthrough
			case actionAppendPartial, actionAppendEscape:
				buf = append(buf, r)
			case actionQuit:
				return
			}
		}
	}()

	return keys, restore, nil
}
