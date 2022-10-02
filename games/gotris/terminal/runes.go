package terminal

type action int

const (
	actionQuit action = iota
	actionAppendAndSend
	actionAppendControl
)

func handleRune(historyLength int, r rune) action {
	switch historyLength {
	case 0:
		switch r {
		case 'q', 3, 4:
			return actionQuit
		case 27:
			return actionAppendControl
		}
	case 1:
		switch r {
		case 91:
			return actionAppendControl
		}
	}
	return actionAppendAndSend
}
