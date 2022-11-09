package fyne

import (
	"context"
	"fmt"

	"fyne.io/fyne/v2"
	"github.com/ubntc/go/games/gotris/input"
)

func (p *Platform) CaptureInput(ctx context.Context) (<-chan *input.Input, func(), error) {
	ch := make(chan *input.Input)
	p.pix.SetOnTypedKey(func(e *fyne.KeyEvent) {
		fmt.Printf("key: %v\n", e)
		// key := Key{*e}
		// ch <- key
	})

	p.pix.SetOnTypedRune(func(r rune) {
		fmt.Printf(" %v ", r)
	})

	cleanup := func() {
		p.pix.SetOnTypedKey(nil)
		p.pix.SetOnTypedRune(nil)
		close(ch)
	}

	return ch, cleanup, nil
}

type Key struct {
	fyne.KeyEvent
}
