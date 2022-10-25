package textscenes

import "strings"

// copied from ans package
const (
	Bold_Red    = "\x1b[1;31m"
	Bold_Green  = "\x1b[1;32m"
	Bold_Yellow = "\x1b[1;33m"
	Red         = "\x1b[31m"
	Reset       = "\x1b[0m"
)

func Colorize(text, pattern, color string) string {
	return strings.ReplaceAll(text, pattern, color+pattern+Reset)
}

func ColorizeBetween(text, startPattern, endPattern, color string) string {
	rows := strings.Split(text, "\n")
	for i, row := range rows {
		l, c, r := "", "", ""
		if lcr := strings.Split(row, startPattern); len(lcr) > 1 {
			l = lcr[0] + startPattern
			c = strings.Join(lcr[1:], startPattern)
		}
		if lcr := strings.Split(c, endPattern); len(lcr) > 1 {
			r = endPattern + lcr[len(lcr)-1]
			c = strings.Join(lcr[:len(lcr)-1], startPattern)
		}

		if l != "" && r != "" {
			rows[i] = l + color + c + Reset + r
		}
	}
	return strings.Join(rows, "\n")
}

func ColorizeFrame(text, verticalFramePattern, contentColor, frameColor string) string {
	text = ColorizeBetween(text, verticalFramePattern, verticalFramePattern, contentColor)
	text = Colorize(text, verticalFramePattern, frameColor)
	return text
}
