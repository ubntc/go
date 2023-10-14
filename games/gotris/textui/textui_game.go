// implements text-based rendering for Gotris
package textui

import (
	"fmt"
	"strings"
	"time"

	geom "github.com/ubntc/go/games/gotris/common/geometry"
	"github.com/ubntc/go/games/gotris/common/platform"
	"github.com/ubntc/go/games/gotris/textui/boxart"
)

func (ui *TextUI) RenderGame(g platform.Game) (rows []string) {
	ch := ui.modeMan.Mode().GetBoxArt()
	var (
		scoreVal = pad(fmt.Sprintf("%d", g.Score), g.PreviewSize.W)
		speedVal = pad(fmt.Sprintf("%d", g.Speed/time.Millisecond), g.PreviewSize.W)

		bw = g.BoardSize.W

		board   = ch.Frame(ui.RenderBoard(g), bw)
		preview = prefix(ch.Space, ui.Title("NEXT", ui.RenderPreview(g)))
		score   = prefix(ch.Space, ui.Title("SCORE", []string{"", scoreVal}))
		level   = prefix(ch.Space, ui.Title("LEVEL", []string{"", speedVal}))
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

	padTop := make([]string, g.BoardPos.H)
	padLeft := strings.Repeat(" ", g.BoardPos.W)

	for i := range res {
		res[i] = padLeft + res[i]
	}

	res = append(padTop, res...)

	return res
}

func (ui *TextUI) RenderPreview(g platform.Game) []string {
	blocks := make(geom.PointMap)
	points := geom.OffsetPointsXY(g.NextTile.Points, 1, 2)
	blocks.SetAll(points, string(g.NextTile.Typ))
	return ui.RenderBlocks(blocks, g.PreviewSize.W, g.PreviewSize.H)
}

func (ui *TextUI) RenderBoard(g platform.Game) (rows []string) {
	blocks := g.Board.Copy()
	geom.MergeTile(g.CurrentTile, blocks)
	return ui.RenderBlocks(blocks, g.BoardSize.W, g.BoardSize.H)
}

func (ui *TextUI) RenderBlocks(blocks geom.PointMap, w, h int) (rows []string) {
	ch := ui.modeMan.Mode().GetBoxArt()
	for y := h - 1; y >= 0; y-- {
		var row []string
		for x := 0; x < w; x++ {
			p := geom.Point{X: x, Y: y}
			name, ok := blocks[p]
			value := ch.BoxMidC
			if ok {
				value = boxart.BlockToString(name, ch.TileCharacters)
			}
			row = append(row, value)
		}
		rows = append(rows, strings.Join(row, ""))
	}
	return
}

func (ui *TextUI) Title(title string, content []string) []string {
	return append([]string{ui.modeMan.Mode().TextToBlock(title)}, content...)
}

func pad(val string, width int) string {
	if len(val) < width {
		val += strings.Repeat(" ", width-len(val))
	}
	return val
}

func prefix(prefix string, rows []string) []string {
	res := make([]string, len(rows))
	for i := range rows {
		res[i] = prefix + rows[i]
	}
	return res
}
