package game

import (
	"github.com/ubntc/go/games/gotris/game/geometry"
	"github.com/ubntc/go/games/gotris/input"
)

type Cmd string

const (
	CmdUnknown     Cmd = "unknown command"
	CmdQuit        Cmd = "quit"
	CmdMoveLeft    Cmd = "move left"
	CmdMoveRight   Cmd = "move right"
	CmdMoveDown    Cmd = "move down"
	CmdMoveUp      Cmd = "move up"
	CmdDrop        Cmd = "drop"
	CmdRotateLeft  Cmd = "rotate left"
	CmdRotateRight Cmd = "rotate right"

	CmdMoveBoardLeft  Cmd = "move board left"
	CmdMoveBoardRight Cmd = "move board right"
	CmdMoveBoardUp    Cmd = "move board up"
	CmdMoveBoardDown  Cmd = "move board down"

	CmdHelp Cmd = "help"
)

func (cmd Cmd) ToDir() geometry.Dir {
	switch cmd {
	case CmdMoveUp:
		return geometry.DirUp
	case CmdMoveDown:
		return geometry.DirDown
	case CmdMoveLeft:
		return geometry.DirLeft
	case CmdMoveRight:
		return geometry.DirRight
	}
	return geometry.DirUnkown
}

func (cmd Cmd) ToSpin() geometry.Spin {
	switch cmd {
	case CmdRotateLeft:
		return geometry.SpinLeft
	case CmdRotateRight:
		return geometry.SpinRight
	}
	return geometry.SpinUnknown
}

func KeyToCmd(key input.Key) (cmd Cmd, ok bool) {
	cmd = CmdUnknown
	defer func() {
		ok = (cmd != CmdUnknown)
	}()

	isMov := (key.Mod() & input.ModMove) != 0
	isAlt := (key.Mod() & input.ModAlt) != 0

	if !isMov {
		key := key.Rune()
		switch key {
		case 'w', 'W':
			// use "WASD up" as additional rotation key to allow one-handed play
			cmd = CmdRotateRight
		case 's', 'S':
			cmd = CmdMoveDown
		case 'a', 'A':
			cmd = CmdMoveLeft
		case 'd', 'D':
			cmd = CmdMoveRight
		case 'z', 'Z':
			cmd = CmdRotateLeft
		case 'x', 'X':
			cmd = CmdRotateRight
		case 'y', 'Y': // Y is next to X German layout
			cmd = CmdRotateLeft
		case 'c', 'C': // setup C + V as alternative keys
			cmd = CmdRotateLeft
		case 'v', 'V': // setup C + V as alternative keys
			cmd = CmdRotateRight
		case 32:
			cmd = CmdDrop
		case 'h', 'H', '?':
			cmd = CmdHelp
		case 229:
			cmd = CmdMoveBoardLeft
		case 8706:
			cmd = CmdMoveBoardRight
		}
		return
	}

	// key is an arrow key movement at this point

	if isAlt {
		switch key.Rune() {
		case 65:
			cmd = CmdMoveBoardUp
		case 66:
			cmd = CmdMoveBoardDown
		case 67, 102:
			cmd = CmdMoveBoardRight
		case 68, 98:
			cmd = CmdMoveBoardLeft
		}
		return
	}

	switch key.Rune() {
	case 65:
		cmd = CmdRotateRight
	case 66:
		cmd = CmdMoveDown
	case 67:
		cmd = CmdMoveRight
	case 68:
		cmd = CmdMoveLeft
	}

	return
}
