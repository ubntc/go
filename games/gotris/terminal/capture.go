package terminal

import (
	"bufio"
	"context"
	"log"
	"os"

	"github.com/pkg/errors"
	xterm "golang.org/x/term"
)

var debug = os.Getenv("DEBUG") != ""

func (t *Terminal) CaptureInput(ctx context.Context) (<-chan []rune, func(), error) {
	runes := make(chan []rune)
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

	sendRunes := func(buf []rune) {
		select {
		case runes <- buf:
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
			close(runes)
		}()
		in := bufio.NewReader(os.Stdin)
		var buf []rune

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
			case actionAppendAndSend:
				sendRunes(append(buf, r))
				buf = nil
			case actionAppendControl:
				buf = append(buf, r)
			case actionQuit:
				return
			}
		}
	}()

	return runes, restore, nil
}
