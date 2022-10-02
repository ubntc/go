package game

type (
	Dir  int
	Spin int
	Cmd  string
)

const (
	DirUnkown Dir = -1
	DirUp     Dir = 0
	DirRight  Dir = 1
	DirDown   Dir = 2
	DirLeft   Dir = 3
)

const (
	SpinUnknown Spin = -1
	SpinLeft    Spin = 0
	SpinRight   Spin = 1
)

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
)

func (cmd Cmd) ToDir() Dir {
	switch cmd {
	case CmdMoveUp:
		return DirUp
	case CmdMoveDown:
		return DirDown
	case CmdMoveLeft:
		return DirLeft
	case CmdMoveRight:
		return DirRight
	}
	return DirUnkown
}

func (cmd Cmd) ToSpin() Spin {
	switch cmd {
	case CmdRotateLeft:
		return SpinLeft
	case CmdRotateRight:
		return SpinRight
	}
	return SpinUnknown
}
