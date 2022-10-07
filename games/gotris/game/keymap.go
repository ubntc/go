package game

func KeyToCmd(keys ...rune) (cmd Cmd, ok bool) {
	cmd = CmdUnknown
	switch len(keys) {
	case 1:
		key := keys[0]
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
		}
	case 2:
	case 3:
		// Arrow Keys [27 91 65:68]
		switch keys[0] {
		case 27:
			switch keys[1] {
			case 91:
				switch keys[2] {
				case 65:
					cmd = CmdRotateRight
				case 66:
					cmd = CmdMoveDown
				case 67:
					cmd = CmdMoveRight
				case 68:
					cmd = CmdMoveLeft
				}
			}
		}
	}

	if cmd != CmdUnknown {
		ok = true
	}
	return
}
