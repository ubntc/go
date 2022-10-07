package rendering

import (
	"fmt"
	"strings"
	"time"

	"github.com/ubntc/go/games/gotris/game"
	"github.com/ubntc/go/games/gotris/game/geometry"
)

func Render(g *game.Game) (rows []string) {
	var (
		scoreVal = Pad(fmt.Sprintf("%d", g.Score), g.PreviewSize.Width)
		speedVal = Pad(fmt.Sprintf("%d", g.Speed/time.Millisecond), g.PreviewSize.Width)

		bw = g.BoardSize.Width
		// pw = g.PreviewSize.Width

		board   = Frame(RenderBoard(g), bw)
		preview = Prefix("  ", Title("ＮＥＸＴ", RenderPreview(g)))
		score   = Prefix("  ", Title("ＳＣＯＲＥ", []string{"", scoreVal}))
		level   = Prefix("  ", Title("ＬＥＶＥＬ", []string{"", speedVal}))
		info    = []string{""}

		res = make([]string, 0, len(board))
	)

	info = append(info, preview...)
	info = append(info, "")
	info = append(info, score...)
	info = append(info, "")
	info = append(info, level...)

	for i, row := range board {
		switch {
		case i < len(info):
			row += info[i]
		default:
			// cannot render info column beyond board height
		}
		res = append(res, row)
	}

	padTop := make([]string, g.BoardPos.Height)
	padLeft := strings.Repeat(" ", g.BoardPos.Width)

	for i := range res {
		res[i] = padLeft + res[i]
	}

	res = append(padTop, res...)

	return res
}

func Title(title string, content []string) []string {
	return append([]string{title}, content...)
}

func Pad(val string, width int) string {
	if len(val) < width {
		val += strings.Repeat(" ", width-len(val))
	}
	return val
}

func Prefix(prefix string, rows []string) []string {
	res := make([]string, len(rows))
	for i := range rows {
		res[i] = prefix + rows[i]
	}
	return res
}

func RenderPreview(g *game.Game) []string {
	blocks := make(geometry.PointMap)
	points := geometry.OffsetPointsXY(g.NextTile.Points(), 1, 2)
	blocks.SetAll(points, game.Blocks[g.NextTile.Typ()])
	return RenderBlocks(blocks, g.PreviewSize.Width, g.PreviewSize.Height)
}

func RenderBoard(g *game.Game) (rows []string) {
	blocks := g.Board.Copy()
	game.MergeTile(g.CurrentTile, blocks)
	return RenderBlocks(blocks, g.BoardSize.Width, g.BoardSize.Height)
}

func RenderBlocks(blocks geometry.PointMap, w, h int) (rows []string) {
	for y := h - 1; y >= 0; y-- {
		var row []string
		for x := 0; x < w; x++ {
			p := geometry.Point{X: x, Y: y}
			color, ok := blocks[p]
			if !ok {
				color = BoxM
			}
			row = append(row, color)
		}
		rows = append(rows, strings.Join(row, ""))
	}
	return
}

var (
	BoxTL, BoxT, BoxTR       = Row("┌一┐") // Top
	BoxL, BoxM, BoxR         = Row("│　│") // Mid
	BoxGndL, BoxGnd, BoxGndR = Row("│￣│") // Ground
	BoxBL, BoxB, BoxBR       = Row("└一┘") // Bottom

	BoxInfoTL, BoxInfoT, BoxInfoTR = Row("　﹏　") // ﹏﹏﹏﹏
	BoxInfoL, BoxInfoM, BoxInfoR   = Row("　　　") // ＴＥＸＴ
	BoxInfoBL, BoxInfoB, BoxInfoBR = Row("　﹋　") // ﹋﹋﹋﹋

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
)

func Row(str string) (l, c, r string) {
	lcr := strings.Split(str, "")
	if len(lcr) != 3 {
		panic("bad box setup")
	}
	return lcr[0], lcr[1], lcr[2]
}

var (
	// blockScore    = "ＳＣＯＲＥ"
	// blockSpeed    = "ＳＰＥＥＤ"
	// blockLevel    = "ＬＥＶＥＬ"
	// blockGameOver = "ＧＡＭＥ　ＯＶＥＲ"
	// blockGotris   = "ＧＯＴＲＩＳ"

	BlockGotrisSmall = []string{
		" ╔═╗╔═╗╔╦╗╦═╗╦╔═╗ ",
		" ║ ╦║ ║ ║ ╠╦╝║╚═╗ ",
		" ╚═╝╚═╝ ╩ ╩╚═╩╚═╝ ",
	}

	TextAbc = "0123456789" + "`" +
		` -+*=/\.,:;!?$%&@#'"<>()[]{}^~_|` +
		`ABCDEFGHIJKLMNOPQRSTUVWXYZ` +
		`abcdefghijklmnopqrstuvwxyz`

	// see: https://www.w3.org/TR/xml-entity-names/023.html
	// and: https://codepoints.net/halfwidth_and_fullwidth_forms

	BlockAbc = strings.Split(`０１２３４５６７８９`+"｀"+
		`　－＋*＝／＼．，：；！？＄％＆＠＃＇＂＜＞（）［］｛｝＾～＿｜`+
		`ＡＢＣＤＥＦＧＨＩＪＫＬＭＮＯＰＱＲＳＴＵＶＷＸＹＺ`+
		`ａｂｃｄｅｆｇｈｉｊｋｌｍｎｏｐｑｒｓｔｕｖｗｘｙｚ`,
		"")
)

func TextToBlock(str string) string {
	res := make([]string, len(str))
	for i, r := range str {
		abcIndex := strings.IndexRune(TextAbc, r)
		res[i] = BlockAbc[abcIndex]
	}
	return strings.Join(res, "")
}

const (
	FullWidthSpace   = "　"
	NormalWidthSpace = " "
	HalfWidthSpace   = "ﾠ"
)

func Frame(rows []string, w int) (frame []string) {
	frame = make([]string, 0, len(rows)+10)

	for _, row := range rows {
		frame = append(frame, BoxL+row+BoxR)
	}

	top := BoxTL + strings.Repeat(BoxT, w) + BoxTR
	gnd := BoxGndL + strings.Repeat(BoxGnd, w) + BoxGndR
	bot := BoxBL + strings.Repeat(BoxB, w) + BoxBR

	frame = append([]string{top}, append(frame, gnd, bot)...)

	return
}
