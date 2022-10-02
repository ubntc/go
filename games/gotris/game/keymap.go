package game

func KeyToCmd(keys ...rune) (cmd Cmd, ok bool) {
	cmd = CmdUnknown
	switch len(keys) {
	case 1:
		key := keys[0]
		switch key {
		case 'd', 'D':
			cmd = CmdMoveDown
		case 'u', 'U':
			cmd = CmdRotateRight
		case 'l', 'L':
			cmd = CmdMoveLeft
		case 'r', 'R':
			cmd = CmdMoveRight
		case 'y', 'Y', 'z', 'Z':
			cmd = CmdRotateLeft
		case 'x', 'X':
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
