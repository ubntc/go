package textscenes

import (
	"strings"
)

type MenuScreen struct {
	RowsTop []string
	RowsBot []string
	PadLeft string
}

const menuKeys = `
╭╭───╮╮           ╭╭───╮╮
││←↕→││---select  ││Any││---confirm
│/‾‾‾\│           │/‾‾‾\│
 ‾‾‾‾‾             ‾‾‾‾‾
`

const confirmOnly = `
╭╭───╮╮
││Any││---confirm
│/‾‾‾\│
 ‾‾‾‾‾
`

var (
	menuKeyRows    = strings.Split(menuKeys, "\n")
	confirmKeyRows = strings.Split(confirmOnly, "\n")
)

const MenuItemsPlaceholder = "MENU_ITEMS"

func NewMenuScreen(source string) MenuScreen {
	screen := MenuScreen{}
	rows := strings.Split(source, "\n")
	i := 0
	for i = range rows {
		if strings.Contains(rows[i], MenuItemsPlaceholder) {
			screen.PadLeft = strings.Split(rows[i], MenuItemsPlaceholder)[0]
			break
		}
	}
	screen.RowsTop = rows[0:i]
	if i < len(rows) {
		screen.RowsBot = rows[i+1:]
	}
	return screen
}

func (screen MenuScreen) menu(names, descriptions []string, current string, keysBlock []string) string {
	var rows []string
	rows = append(rows, screen.RowsTop...)

	for i, name := range names {
		pad := screen.PadLeft
		if name == current {
			pad += "→ "
		} else {
			pad += "  "
		}
		desc := ""
		if len(descriptions) > i {
			desc = ": " + descriptions[i]
		}
		rows = append(rows, pad+name+desc)
	}

	for _, row := range keysBlock {
		rows = append(rows, screen.PadLeft+row)
	}

	rows = append(rows, screen.RowsBot...)

	return strings.Join(rows, "\n")
}

func (screen MenuScreen) Menu(names, descriptions []string, current string) string {
	return screen.menu(names, descriptions, current, nil)
}

func Screen(screen string) string {
	return NewMenuScreen(screen).menu(nil, nil, "", confirmKeyRows)
}
