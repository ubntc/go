package controls

import (
	"github.com/ubntc/go/games/gotris/input"
)

func InputToCmd(in input.Input) (c Cmd, arg string) {
	if in.IsEmpty() {
		return
	}

	switch {
	case in.IsAlt():
		switch in.Key() {
		case input.KeyUp:
			c = MoveBoardUp
		case input.KeyDown:
			c = MoveBoardDown
		case input.KeyRight:
			c = MoveBoardRight
		case input.KeyLeft:
			c = MoveBoardLeft
		}
	default:
		switch in.Key() {
		case input.KeyUp:
			// use "up" as additional rotation key to allow one-handed play
			c = RotateRight
		case input.KeyDown:
			c = MoveDown
		case input.KeyLeft:
			c = MoveLeft
		case input.KeyRight:
			c = MoveRight
		case input.KeyButton1:
			c = RotateLeft
		case input.KeyButton2:
			c = RotateRight
		case input.KeyButton3:
			c = Drop
		case input.KeyHelp:
			c = Help
		case input.KeyMenu:
			c = Options
		case input.KeyQuit:
			c = Quit
		case input.KeyNumber:
			c = SelectMode
			arg = string(in.Rune())
		}
	}

	return
}

func KeyToMenuCmd(in input.Input) (cmd Cmd, arg string) {
	if in.IsEmpty() {
		return
	}

	switch in.Key() {
	case input.KeyQuit:
		cmd = Quit
	case input.KeyUp:
		cmd = MenuUp
	case input.KeyDown:
		cmd = MenuDown
	case input.KeyLeft:
		cmd = MenuLeft
	case input.KeyRight:
		cmd = MenuRight
	case input.KeyHelp:
		cmd = Help
	case input.KeyEnter, input.KeyButton3:
		cmd = MenuSelect
	}

	return
}
