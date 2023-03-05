// generic ascii/unicode/ansi arts interfaces and types
package arts

import (
	"fmt"
	"strings"
)

type FrameArt interface {
	Art() *Characters
	TextToBlock(string) string
}

type Characters struct {
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

	TextAbc  string
	BlockAbc []string
}

func (fa *Characters) Examples() (blocks string) {
	for _, c := range fa.TileCharacters {
		blocks += c
	}
	return
}

func (fa *Characters) Info() (text string) {
	return fmt.Sprintf("%s: %s", fa.Name, fa.Desc)
}

func (fa *Characters) TextToBlock(text string) (block string) {
	return text
}

func Row(lcr ...string) (l, c, r string) {
	if len(lcr) == 1 {
		lcr = strings.Split(lcr[0], "")
	}

	n := len(lcr)
	if n < 3 {
		panic(fmt.Sprintf("bad box strings: %s", lcr))
	}
	return lcr[0], strings.Join(lcr[1:n-1], ""), lcr[n-1]
}

func BlockToString(name string, chars map[string]string) string {
	if s, ok := chars[name]; ok {
		return s
	}
	return "?"
}
