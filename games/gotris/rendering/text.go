package rendering

import (
	"fmt"
	"strings"
	"time"

	"github.com/ubntc/go/games/gotris/game"
)

func Render(g *game.Game) (rows []string) {
	var (
		scoreVal = Pad(fmt.Sprintf("%d", g.Score), g.PreviewSize.Width)
		speedVal = Pad(fmt.Sprintf("%d ms", g.Speed/time.Millisecond), g.PreviewSize.Width)

		bw = g.BoardSize.Width
		// pw = g.PreviewSize.Width

		board   = Frame(RenderBoard(g), bw)
		preview = Prefix("  ", Title("NEXT", RenderPreview(g)))
		score   = Prefix("  ", Title("SCORE", []string{"", scoreVal}))
		speed   = Prefix("  ", Title("Speed", []string{"", speedVal}))
		info    = []string{""}

		res = make([]string, 0, len(board))
	)

	info = append(info, preview...)
	info = append(info, "")
	info = append(info, score...)
	info = append(info, "")
	info = append(info, speed...)

	for i, row := range board {
		switch {
		case i < len(info):
			row += info[i]
		default:
			// cannot render info column beyond board height
		}
		res = append(res, row)
	}
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
	blocks := make(game.PointMap)
	points := game.OffsetPointsXY(g.NextTile.Points(), 1, 2)
	game.MergePoints(points, game.Blocks[g.NextTile.Typ()], blocks)
	return RenderBlocks(blocks, g.PreviewSize.Width, g.PreviewSize.Height)
}

func RenderBoard(g *game.Game) (rows []string) {
	blocks := g.Board.Copy()
	game.MergeTile(g.CurrentTile, blocks)
	return RenderBlocks(blocks, g.BoardSize.Width, g.BoardSize.Height)
}

func RenderBlocks(blocks game.PointMap, w, h int) (rows []string) {
	for y := h - 1; y >= 0; y-- {
		var row []string
		for x := 0; x < w; x++ {
			p := game.Point{X: x, Y: y}
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
	BoxTL = ""
	BoxT  = ""
	BoxTR = ""
	BoxL  = ""
	BoxM  = ""
	BoxR  = ""
	BoxBL = ""
	BoxB  = ""
	BoxBR = ""
)

func SetFrameCharacters(box []string) {
	BoxTL = strings.Split(box[0], "")[0]
	BoxT = strings.Split(box[0], "")[1]
	BoxTR = strings.Split(box[0], "")[2]
	BoxL = strings.Split(box[1], "")[0]
	BoxM = strings.Split(box[1], "")[1]
	BoxR = strings.Split(box[1], "")[2]
	BoxBL = strings.Split(box[2], "")[0]
	BoxB = strings.Split(box[2], "")[1]
	BoxBR = strings.Split(box[2], "")[2]
}

func init() {
	SetFrameCharacters([]string{
		// see: https://www.w3.org/TR/xml-entity-names/023.html
		"⎛﹋⎞", // "╒═╕",
		"│　│", // "│ │",
		"⎝﹏⎠", // "╘═╛",
	})
}

// DoubleWidthSpace defines the characters to draw the game background.
const DoubleWidthSpace = "　"
const SingleWidthSpace = " "

func Frame(rows []string, w int) (frame []string) {
	frame = make([]string, 0, len(rows)+10)

	for _, row := range rows {
		frame = append(frame, BoxL+row+BoxR)
	}

	top := BoxTL + strings.Repeat(BoxT, w) + BoxTR
	mid := BoxL + strings.Repeat(BoxT, w) + BoxR
	bot := BoxBL + strings.Repeat(BoxB, w) + BoxBR

	frame = append([]string{top}, append(frame, mid, bot)...)

	return
}
