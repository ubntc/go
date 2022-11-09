package input

// NewFromRune create a game-focussed input event by mapping
// common key like WASD to movement keys. The original rune will
// be preserved and can be accessed viaInput.Rune().
// For some obscure runes, the rune will be set to 0.
func NewFromRune(r rune, flags Flag) *Input {
	key := KeyNone

	switch {
	case flags.IsMovement():
		switch r {
		case 65:
			key = KeyUp
		case 66:
			key = KeyDown
		case 67, 102:
			key = KeyRight
		case 68, 98:
			key = KeyLeft
		}
		r = 0
	default:
		switch r {
		case 'w', 'W':
			key = KeyUp
		case 's', 'S':
			key = KeyDown
		case 'a', 'A':
			key = KeyLeft
		case 'd', 'D':
			key = KeyRight
		case 'z', 'Z', 'c', 'C', 'y', 'Y':
			key = KeyButton1
		case 'x', 'X', 'v', 'V':
			key = KeyButton2
		case ' ':
			key = KeyButton3
		case 'q', 'Q':
			key = KeyQuit
		case 13:
			key = KeyEnter
		case 'h', 'H', '?':
			key = KeyHelp
		case 'o', 'O', ',', 'm', 'M':
			key = KeyMenu
		case 229: // Alt Left variant
			key = KeyLeft
			flags |= FlagAlt | FlagMove
			r = 0
		case 8706: // Alt Right variant
			key = KeyRight
			flags |= FlagAlt | FlagMove
			r = 0
		case '1', '2', '3', '4', '5', '6', '7', '8', '9', '0':
			key = KeyNumber
		default:
			key = KeyRune
		}
	}

	return &Input{key, flags, r}
}
