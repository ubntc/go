// generic ascii/unicode/ansi arts interfaces and types
package boxart

import (
	"fmt"
	"strings"
)

type BoxArt struct {
	Name string
	Desc string

	TileCharacters map[string]string
	BlockGotris    []string

	Space string

	BoxTopL, BoxTopC, BoxTopR string
	BoxMidL, BoxMidC, BoxMidR string
	BoxGndL, BoxGndC, BoxGndR string
	BoxBotL, BoxBotC, BoxBotR string

	BoxInfoTL, BoxInfoT, BoxInfoTR string
	BoxInfoML, BoxInfoC, BoxInfoMR string
	BoxInfoBL, BoxInfoB, BoxInfoBR string

	TextAbc  string   // TextAbc defines all supported text characters
	BlockAbc []string // BlockAbc defines a rune block for each text character in TextAbc
}

// GetBoxArt returns the BoxArt object itself. Use this method to access the concrete BoxArt
// in embedded types. Also see mode.Frameboxart.
func (a *BoxArt) GetBoxArt() *BoxArt {
	return a
}

func (a *BoxArt) GetName() string {
	if a.Name == "" {
		return "unnamed BoxArt"
	}
	return a.Name
}

func (a *BoxArt) GetDesc() string {
	if a.Desc == "" {
		return "BoxArt:" + a.GetName()
	}
	return a.Desc
}

func (a *BoxArt) Examples() (blocks string) {
	for _, c := range a.TileCharacters {
		blocks += c
	}
	return
}

func (a *BoxArt) Info() (text string) {
	return fmt.Sprintf("%s: %s", a.Name, a.Desc)
}

// TextToBlock looks up the all chars of the given str in the TextAbc
// and returns the corresponding rune block from BlockAbc.
//
// This is the default implementation of a text-to-boxart mapping function
// that processes the input rune by rune. This function can be overwritten
// in derived types. See doublewidth.doublewidth as an example.
func (a *BoxArt) TextToBlock(str string) string {
	res := make([]string, len(str))
	for i, r := range str {
		abcIndex := strings.IndexRune(a.TextAbc, r)
		res[i] = a.BlockAbc[abcIndex]
	}
	return strings.Join(res, "")
}

// Frame returns the given rows wrapped in the BoxArt's frame runes.
func (a *BoxArt) Frame(rows []string, w int) (frame []string) {
	frame = make([]string, 0, len(rows)+10)

	for _, row := range rows {
		frame = append(frame, a.BoxMidL+row+a.BoxMidR)
	}

	top := a.BoxTopL + strings.Repeat(a.BoxTopC, w) + a.BoxTopR
	gnd := a.BoxGndL + strings.Repeat(a.BoxGndC, w) + a.BoxGndR
	bot := a.BoxBotL + strings.Repeat(a.BoxBotC, w) + a.BoxBotR

	frame = append([]string{top}, append(frame, gnd, bot)...)

	return
}
