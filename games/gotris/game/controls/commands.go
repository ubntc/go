package controls

import (
	"github.com/ubntc/go/games/gotris/game/geometry"
	"github.com/ubntc/go/games/gotris/input"
)

type Cmd string

const (
	Empty Cmd = ""

	Quit        Cmd = "quit"
	MoveLeft    Cmd = "move left"
	MoveRight   Cmd = "move right"
	MoveDown    Cmd = "move down"
	MoveUp      Cmd = "move up"
	Drop        Cmd = "drop"
	RotateLeft  Cmd = "rotate left"
	RotateRight Cmd = "rotate right"

	MoveBoardLeft  Cmd = "move board left"
	MoveBoardRight Cmd = "move board right"
	MoveBoardUp    Cmd = "move board up"
	MoveBoardDown  Cmd = "move board down"

	SelectMode Cmd = "select rendering mode in-game"

	Help    Cmd = "help"
	Options Cmd = "options"
)

func (c Cmd) ToDir() geometry.Dir {
	switch c {
	case MoveUp:
		return geometry.DirUp
	case MoveDown:
		return geometry.DirDown
	case MoveLeft:
		return geometry.DirLeft
	case MoveRight:
		return geometry.DirRight
	}
	return geometry.DirUnkown
}

func (c Cmd) ToSpin() geometry.Spin {
	switch c {
	case RotateLeft:
		return geometry.SpinLeft
	case RotateRight:
		return geometry.SpinRight
	}
	return geometry.SpinUnknown
}

func KeyToCmd(key input.Key) (c Cmd, arg string) {
	if key == nil {
		return
	}

	if !input.IsMovement(key) {
		key := key.Rune()
		switch key {
		case 'w', 'W':
			// use "WASD up" as additional rotation key to allow one-handed play
			c = RotateRight
		case 's', 'S':
			c = MoveDown
		case 'a', 'A':
			c = MoveLeft
		case 'd', 'D':
			c = MoveRight
		case 'z', 'Z':
			c = RotateLeft
		case 'x', 'X':
			c = RotateRight
		case 'y', 'Y': // Y is next to X German layout
			c = RotateLeft
		case 'c', 'C': // setup C + V as alternative keys
			c = RotateLeft
		case 'v', 'V': // setup C + V as alternative keys
			c = RotateRight
		case ' ':
			c = Drop
		case 'h', 'H', '?':
			c = Help
		case 'o', 'O', ',':
			c = Options
		case 229:
			c = MoveBoardLeft
		case 8706:
			c = MoveBoardRight
		case '1', '2', '3', '4', '5', '6', '7', '8', '9', '0':
			c = SelectMode
			arg = string(key)
		}
		return
	}

	// key is an arrow key movement at this point
	if input.IsAlt(key) {

		switch key.Rune() {
		case 65:
			c = MoveBoardUp
		case 66:
			c = MoveBoardDown
		case 67, 102:
			c = MoveBoardRight
		case 68, 98:
			c = MoveBoardLeft
		}
		return
	}

	switch key.Rune() {
	case 65:
		c = RotateRight
	case 66:
		c = MoveDown
	case 67:
		c = MoveRight
	case 68:
		c = MoveLeft
	}

	return
}

const (
	MenuUp     Cmd = "menu up"
	MenuDown   Cmd = "menu down"
	MenuLeft   Cmd = "menu left"
	MenuRight  Cmd = "menu right"
	MenuSelect Cmd = "menu select"
)

func KeyToMenuCmd(key input.Key) (cmd Cmd, arg string) {
	if key == nil {
		return
	}

	if !input.IsMovement(key) {
		switch key.Rune() {
		case 'w', 'W':
			// use "WASD up" as additional rotation key to allow one-handed play
			cmd = MenuUp
		case 's', 'S':
			cmd = MenuDown
		case 'a', 'A':
			cmd = MenuLeft
		case 'd', 'D':
			cmd = MenuRight
		case 'h', 'H', '?':
			cmd = Help
		case 13, ' ':
			cmd = MenuSelect
		}
		return
	} else {
		switch key.Rune() {
		case 65:
			cmd = MenuUp
		case 66:
			cmd = MenuDown
		case 67:
			cmd = MenuRight
		case 68:
			cmd = MenuLeft
		}
	}
	return
}

func HandleOptionsCmd(command Cmd, numOptions, idx int) (newIndex int, done, ok bool) {
	switch command {
	case Empty, MenuSelect:
		done = true
	case MenuDown, MenuRight:
		idx = (idx + 1) % numOptions
	case MenuUp, MenuLeft:
		idx = (idx + numOptions - 1) % numOptions
	default:
		// do not handle and other commands indicated using (ok = false)
		return -1, false, false
	}
	// command was handled (ok = true)
	return idx, done, true
}
