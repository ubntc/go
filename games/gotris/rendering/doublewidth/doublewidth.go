package doublewidth

import (
	"strings"

	"github.com/ubntc/go/games/gotris/rendering/arts"
)

type doublewidth struct{ arts.Characters }

const (
	canvas_block = arts.Block_BigFrame
	bg_canvas    = arts.BG_D_Gray
	canvasSpace  = bg_canvas + arts.FG_D_Gray + canvas_block + arts.Reset

	text  = arts.Block_SmallDoubleFrame
	block = arts.Bold + text
	bg    = bg_canvas
	end   = arts.Reset
	start = bg
)

func New() arts.FrameArt {
	c := &doublewidth{}
	c.Space = "  "
	c.Name = "Double-Width"
	c.Desc = "block-width: 2ch, block-char: '" + text + "', frames: ascii"

	c.TileCharacters = map[string]string{
		"B": start + arts.FG_Yellow + block + end, // std: yellow | "🟨",
		"I": start + arts.FG_Cyan__ + block + end, // std: cyan   | "🟧",
		"L": start + arts.FG_Orange + block + end, // std: orange | "🟫",
		"J": start + arts.FG_Blue__ + block + end, // std: blue   | "🟦",
		"T": start + arts.FG_Purple + block + end, // std: purple | "🟪",
		"S": start + arts.FG_Green_ + block + end, // std: green  | "🟩",
		"Z": start + arts.FG_Red___ + block + end, // std: red    | "🟥",
	}

	c.BlockGotris = []string{
		" ╔═╗╔═╗╔╦╗╦═╗╦╔═╗ ",
		" ║ ╦║ ║ ║ ╠╦╝║╚═╗ ",
		" ╚═╝╚═╝ ╩ ╩╚═╩╚═╝ ",
	}

	c.BoxTopL, c.BoxTopC, c.BoxTopR = arts.Row("╭──╮") // Top     ┌──┐ ╭──╮ ╒══╕ ╒══╕ ╔══╗
	c.BoxMidL, c.BoxMidC, c.BoxMidR = arts.Row("│  │") // Mid     │  │ │  │ │  │ │  │ ║  ║
	c.BoxGndL, c.BoxGndC, c.BoxGndR = arts.Row("├──┤") // Ground  ├──┤ ├──┤ ├──┤ ╞══╡ ╠══╣
	c.BoxBotL, c.BoxBotC, c.BoxBotR = arts.Row("╰──╯") // Bottom  └──┘ ╰──╯ ╘══╛ ╘══╛ ╚══╝

	c.BoxInfoTL, c.BoxInfoT, c.BoxInfoTR = arts.Row(" _ ")
	c.BoxInfoML, c.BoxInfoC, c.BoxInfoMR = arts.Row("   ")
	c.BoxInfoBL, c.BoxInfoB, c.BoxInfoBR = arts.Row(" ‾ ")

	c.BoxMidC = canvasSpace

	c.TextAbc = "0123456789" + "`" +
		` -+*=/\.,:;!?$%&@#'"<>()[]{}^~_|` +
		`ABCDEFGHIJKLMNOPQRSTUVWXYZ` +
		`abcdefghijklmnopqrstuvwxyz`

	c.BlockAbc = strings.Split(c.TextAbc, "")
	return c
}

func (c *doublewidth) Art() *arts.Characters { return &c.Characters }

func (c *doublewidth) TextToBlock(str string) string {
	res := make([]string, len(str))
	for i, r := range str {
		abcIndex := strings.IndexRune(c.TextAbc, r)
		res[i] = c.BlockAbc[abcIndex]
	}
	return strings.Join(res, "")
}
