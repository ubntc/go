package input

type Key string

// Flag indicates a key modifier or key group
type Flag int

const (
	FlagNone  Flag = 0
	FlagShift Flag = 1
	FlagAlt   Flag = 2
	FlagCtrl  Flag = 4
	FlagMove  Flag = 8
)

func (f Flag) IsMovement() bool { return (f & FlagMove) != 0 }
func (f Flag) IsAlt() bool      { return (f & FlagAlt) != 0 }

const (
	KeyNone    Key = ""
	KeyUp      Key = "up"
	KeyDown    Key = "down"
	KeyLeft    Key = "left"
	KeyRight   Key = "right"
	KeyEnter   Key = "enter"
	KeyQuit    Key = "quit"
	KeyHelp    Key = "help"
	KeyButton1 Key = "button 1"
	KeyButton2 Key = "button 2"
	KeyButton3 Key = "button 3"
	KeyMenu    Key = "menu"
	KeyNumber  Key = "number"
	KeyRune    Key = "rune"
)

// Keys is a slice of all supported keys.
// It should be used by a game platform to test that all keys are mapped.
var Keys = []Key{
	KeyNone,
	KeyUp,
	KeyDown,
	KeyLeft,
	KeyRight,
	KeyEnter,
	KeyButton1,
	KeyButton2,
	KeyQuit,
}
