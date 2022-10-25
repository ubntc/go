package fullwidth

import (
	"strings"

	"github.com/ubntc/go/games/gotris/rendering/text/arts"
)

const Name = "Full-Width"

type fullwidth struct{ arts.Characters }

const space = arts.FullWidthSpace

func New() arts.FrameArt {
	c := &fullwidth{}
	c.Space = space
	c.Name = Name
	c.Desc = "block-width: 1cjk, block-char: '🟪', frames: ascii/unicode"

	c.TileCharacters = map[string]string{
		"B": "🟫",
		"I": "🟨",
		"L": "🟥",
		"J": "🟧",
		"T": "🟩",
		"S": "🟦",
		"Z": "🟪",
	}

	c.BoxTopL, c.BoxTopC, c.BoxTopR = arts.Row("┌一┐") // Top
	c.BoxMidL, c.BoxMidC, c.BoxMidR = arts.Row("│・│") // Mid
	c.BoxGndL, c.BoxGndC, c.BoxGndR = arts.Row("│￣│") // Ground
	c.BoxBotL, c.BoxBotC, c.BoxBotR = arts.Row("└一┘") // Bottom

	c.BoxInfoTL, c.BoxInfoT, c.BoxInfoTR = arts.Row("　﹏　") // ﹏﹏﹏﹏
	c.BoxInfoML, c.BoxInfoC, c.BoxInfoMR = arts.Row("　　　") // ＴＥＸＴ
	c.BoxInfoBL, c.BoxInfoB, c.BoxInfoBR = arts.Row("　﹋　") // ﹋﹋﹋﹋

	// ・一一一一一・ . 　＿＿＿＿＿　 . 　＿＿＿＿＿
	// ｜　　　　　｜ . ｜　　　　　｜ . ｜　　　　　｜
	// ｜　　🟩　　｜ . ｜　　🟩　　｜ . ｜　　🟩　　｜
	// ｜　🟩🟩🟩　｜ . ｜　🟩🟩🟩　｜ . ｜　🟩🟩🟩　｜
	// ｜　　　　　｜ . ｜　　　　　｜ . ｜　　　　　｜
	// ｜￣￣￣￣￣｜ . ｜￣￣￣￣￣｜ . ｜⬛️⬛️⬛️⬛️⬛️｜
	// ・一一一一一・ . 　￣￣￣￣￣　 . 　￣￣￣￣￣
	//
	// ╒＝＝＝＝╕ . ╒－－－－╕ . ┌一一一一┐
	// │　　　　│ . │　　　　│ . │　　　　│
	// │　🟩🟩　│ . │　🟩🟩　│ . │　🟩🟩　│
	// │　　　　│ . │　　　　│ . │　　　　│
	// │￣￣￣￣│ . │￣￣￣￣│ . │￣￣￣￣│
	// ╘＝＝＝＝╛ . ╘－－－－╛ . └一一一一┘

	c.TextAbc = "0123456789" + "`" +
		` -+*=/\.,:;!?$%&@#'"<>()[]{}^~_|` +
		`ABCDEFGHIJKLMNOPQRSTUVWXYZ` +
		`abcdefghijklmnopqrstuvwxyz`

	// see: https://www.w3.org/TR/xml-entity-names/023.html
	// and: https://codepoints.net/halfwidth_and_fullwidth_forms

	c.BlockAbc = strings.Split(`０１２３４５６７８９`+"｀"+
		`　－＋*＝／＼．，：；！？＄％＆＠＃＇＂＜＞（）［］｛｝＾～＿｜`+
		`ＡＢＣＤＥＦＧＨＩＪＫＬＭＮＯＰＱＲＳＴＵＶＷＸＹＺ`+
		`ａｂｃｄｅｆｇｈｉｊｋｌｍｎｏｐｑｒｓｔｕｖｗｘｙｚ`,
		"")

	return c
}

func (c *fullwidth) Art() *arts.Characters { return &c.Characters }

func (c *fullwidth) TextToBlock(str string) string {
	res := make([]string, len(str))
	for i, r := range str {
		abcIndex := strings.IndexRune(c.TextAbc, r)
		res[i] = c.BlockAbc[abcIndex]
	}
	return strings.Join(res, "")
}
