package controls

import (
	"github.com/ubntc/go/games/gotris/common/geometry"
	"github.com/ubntc/go/games/gotris/common/options"
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

	MenuUp     Cmd = "menu up"
	MenuDown   Cmd = "menu down"
	MenuLeft   Cmd = "menu left"
	MenuRight  Cmd = "menu right"
	MenuSelect Cmd = "menu select"
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

type HandleResult int

const (
	HandleResultNotHandled HandleResult = iota
	HandleResultSelectionFinished
	HandleResultSelectionChanged
)

func HandleOptionsCmd(command Cmd, options options.Options) HandleResult {
	switch command {
	case Empty, MenuSelect:
		return HandleResultSelectionFinished
	case MenuDown, MenuRight:
		options.Inc()
		return HandleResultSelectionChanged
	case MenuUp, MenuLeft:
		options.Dec()
		return HandleResultSelectionChanged
	default:
		return HandleResultNotHandled
	}
}
