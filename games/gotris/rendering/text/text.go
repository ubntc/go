package text

import (
	"fmt"
	"strings"
	"time"

	"github.com/ubntc/go/games/gotris/game"
	"github.com/ubntc/go/games/gotris/game/geometry"
	"github.com/ubntc/go/games/gotris/game/tiles"
	"github.com/ubntc/go/games/gotris/rendering/text/arts"
	"github.com/ubntc/go/games/gotris/rendering/text/modes"
)

var modeMan = modes.NewModeManager()

func art() arts.FrameArt          { return modeMan.Mode() }
func ch() *arts.Characters        { return modeMan.Mode().Art() }
func ModeMan() *modes.ModeManager { return modeMan }

func Render(g *game.Game) (rows []string) {
	ch := ch()
	var (
		scoreVal = Pad(fmt.Sprintf("%d", g.Score), g.PreviewSize.W)
		speedVal = Pad(fmt.Sprintf("%d", g.Speed/time.Millisecond), g.PreviewSize.W)

		bw = g.BoardSize.W
		// pw = g.PreviewSize.Width

		board   = Frame(RenderBoard(g), bw)
		preview = Prefix(ch.Space, Title("NEXT", RenderPreview(g)))
		score   = Prefix(ch.Space, Title("SCORE", []string{"", scoreVal}))
		level   = Prefix(ch.Space, Title("LEVEL", []string{"", speedVal}))
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

	pos := g.BoardPos
	padTop := make([]string, pos.H)
	padLeft := strings.Repeat(" ", pos.W)

	for i := range res {
		res[i] = padLeft + res[i]
	}

	res = append(padTop, res...)

	return res
}

func Title(title string, content []string) []string {
	return append([]string{art().TextToBlock(title)}, content...)
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
	blocks.SetAll(points, string(g.NextTile.Typ()))
	return RenderBlocks(blocks, g.PreviewSize.W, g.PreviewSize.H)
}

func RenderBoard(g *game.Game) (rows []string) {
	blocks := g.Board.Copy()
	tiles.MergeTile(g.CurrentTile, blocks)
	return RenderBlocks(blocks, g.BoardSize.W, g.BoardSize.H)
}

func RenderBlocks(blocks geometry.PointMap, w, h int) (rows []string) {
	ch := art().Art()
	for y := h - 1; y >= 0; y-- {
		var row []string
		for x := 0; x < w; x++ {
			p := geometry.Point{X: x, Y: y}
			name, ok := blocks[p]
			value := ch.BoxMidC
			if ok {
				value = arts.BlockToString(name, ch.TileCharacters)
			}
			row = append(row, value)
		}
		rows = append(rows, strings.Join(row, ""))
	}
	return
}

func Frame(rows []string, w int) (frame []string) {
	ch := art().Art()
	frame = make([]string, 0, len(rows)+10)

	for _, row := range rows {
		frame = append(frame, ch.BoxMidL+row+ch.BoxMidR)
	}

	top := ch.BoxTopL + strings.Repeat(ch.BoxTopC, w) + ch.BoxTopR
	gnd := ch.BoxGndL + strings.Repeat(ch.BoxGndC, w) + ch.BoxGndR
	bot := ch.BoxBotL + strings.Repeat(ch.BoxBotC, w) + ch.BoxBotR

	frame = append([]string{top}, append(frame, gnd, bot)...)

	return
}
