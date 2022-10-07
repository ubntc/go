package terminal

type action int

const (
	actionQuit action = iota
	actionAppendAndSend
	actionAppendEscape
	actionAppendPartial
	actionAppendShift
	actionSendWithAltAsMovement
	actionAppendAlt
	actionAppendAltShift
	actionAppendCtrl
	actionAppendCtrlShift
	actionAppendMovement
)

func handleRune(historyLength int, r rune) action {
	switch historyLength {
	case 0:
		switch r {
		case 'q', 3, 4:
			return actionQuit
		case 27:
			// Escape charcter (aka. 033, 0x1B), part 1/?
			return actionAppendEscape
		}
	case 1:
		switch r {
		case 91:
			// Escape part 2/?, used by most arrow key combinations
			return actionAppendMovement
		case 98, 102:
			// Escape part 2/2, without further keys pending, used only for ALT + ←→
			return actionSendWithAltAsMovement
		case 27:
			// iTerm2 send this on ALT + ←→
			return actionAppendAlt
		}
	case 2:
		switch r {
		case 49, 91:
			// Escape part 3/?, modfier part 1/2
			return actionAppendMovement
		}
	case 3:
		switch r {
		case 59:
			// Escape part 4/5, modifer part 2/2
			return actionAppendPartial
		}
	case 4:
		switch r {
		case 50:
			return actionAppendShift
		case 51:
			return actionAppendAlt
		case 52:
			return actionAppendAltShift
		case 53:
			return actionAppendCtrl
		case 54:
			return actionAppendCtrlShift
		}
	}
	return actionAppendAndSend
}
