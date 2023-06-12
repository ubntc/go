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

func (t *Terminal) CaptureInput(ctx context.Context) (<-chan input.Input, func(), error) {
	keys := make(chan input.Input, 10)
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

	sendKey := func(in *input.Input) {
		select {
		case keys <- *in:
		default:
			// ignore new input if prev. input is stalled
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
		var mod input.Flag

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
				mod |= input.FlagAlt | input.FlagMove
				fallthrough
			case actionAppendAndSend:
				sendKey(input.NewFromRune(r, mod))
				mod = 0
				buf = nil
			case actionAppendAlt:
				mod |= input.FlagAlt
			case actionAppendCtrl:
				mod |= input.FlagCtrl
			case actionAppendShift:
				mod |= input.FlagShift
			case actionAppendAltShift:
				mod |= input.FlagAlt | input.FlagShift
			case actionAppendCtrlShift:
				mod |= input.FlagCtrl | input.FlagShift
			case actionAppendMovement:
				mod |= input.FlagMove
				fallthrough
			case actionAppendPartial, actionAppendEscape:
				buf = append(buf, r)
			case actionQuit:
				// TODO: allow Q as rune and not just instant quit
				return
			}
		}
	}()

	return keys, restore, nil
}
