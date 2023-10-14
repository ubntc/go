package fullwidth

import (
	"strings"

	"github.com/ubntc/go/games/gotris/textui/boxart"
)

const Name = "Full-Width"

type fullwidth struct {
	boxart.BoxArt // embedded BoxArt type gives access to BoxArt methods
}

func New() *fullwidth {
	c := &fullwidth{}
	c.Space = boxart.FullWidthSpace
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

	c.BoxTopL, c.BoxTopC, c.BoxTopR = boxart.Row("┌一┐") // Top
	c.BoxMidL, c.BoxMidC, c.BoxMidR = boxart.Row("│・│") // Mid
	c.BoxGndL, c.BoxGndC, c.BoxGndR = boxart.Row("│￣│") // Ground
	c.BoxBotL, c.BoxBotC, c.BoxBotR = boxart.Row("└一┘") // Bottom

	c.BoxInfoTL, c.BoxInfoT, c.BoxInfoTR = boxart.Row("　﹏　") // ﹏﹏﹏﹏
	c.BoxInfoML, c.BoxInfoC, c.BoxInfoMR = boxart.Row("　　　") // ＴＥＸＴ
	c.BoxInfoBL, c.BoxInfoB, c.BoxInfoBR = boxart.Row("　﹋　") // ﹋﹋﹋﹋

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
