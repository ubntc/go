package fyne

import (
	"context"
	"fmt"

	"fyne.io/fyne/v2"
	"github.com/ubntc/go/games/gotris/input"
)

func (p *Platform) CaptureInput(ctx context.Context) (<-chan input.Input, func(), error) {
	ch := make(chan input.Input)
	p.pix.SetOnTypedKey(func(e *fyne.KeyEvent) {
		fmt.Println("fyne key:", e)
		key := InputKeyFromKeyName(e)
		ch <- *input.New(key, 0)
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

func InputKeyFromKeyName(k *fyne.KeyEvent) input.Key {
	switch k.Name {
	case fyne.KeyLeft:
		return input.KeyLeft
	case fyne.KeyRight:
		return input.KeyRight
	case fyne.KeyUp:
		return input.KeyUp
	case fyne.KeyDown:
		return input.KeyDown
	case fyne.KeyY, fyne.KeyC, fyne.KeyZ:
		return input.KeyButton1
	case fyne.KeyX, fyne.KeyV:
		return input.KeyButton2
	case fyne.KeySpace:
		return input.KeyButton3
	case fyne.KeyEnter, fyne.KeyReturn:
		return input.KeyEnter
	case fyne.KeyQ:
		return input.KeyQuit
	case fyne.KeyH:
		return input.KeyHelp
	case fyne.KeyM, fyne.KeyO, fyne.KeyComma:
		return input.KeyMenu
	}
	return input.KeyNone
}
