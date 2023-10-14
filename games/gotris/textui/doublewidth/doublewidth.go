package doublewidth

import (
	"strings"

	"github.com/ubntc/go/games/gotris/textui/boxart"
)

const Name = "Double-Width"

type doublewidth struct {
	boxart.BoxArt // embedded BoxArt type gives access to BoxArt methods
}

const (
	bg    = boxart.BgDkGry                // default BG color used for all canvas drawing
	text  = boxart.Block_SmallDoubleFrame // unicode char that represents a single block
	block = boxart.Bold + text            // styled block

	end   = boxart.Reset // at the end of each block reset and modifications
	start = bg           // all blocks start with settng a background color

	// canvas overrides and spaces dervied from the inline box generation below
	canvas = start + boxart.DkGry + boxart.Block_BigFrame + end
)

func New() *doublewidth {
	c := &doublewidth{}
	c.Space = "  "
	c.Name = Name
	c.Desc = "block-width: 2ch, block-char: '" + text + "', frames: ascii"

	c.TileCharacters = map[string]string{
		"B": start + boxart.Yel + block + end, // std: yellow
		"I": start + boxart.Cyn + block + end, // std: cyan
		"L": start + boxart.Ora + block + end, // std: orange
		"J": start + boxart.Blu + block + end, // std: blue
		"T": start + boxart.Pur + block + end, // std: purple
		"S": start + boxart.Grn + block + end, // std: green
		"Z": start + boxart.Red + block + end, // std: red
	}

	c.BlockGotris = []string{
		" ╔═╗╔═╗╔╦╗╦═╗╦╔═╗ ",
		" ║ ╦║ ║ ║ ╠╦╝║╚═╗ ",
		" ╚═╝╚═╝ ╩ ╩╚═╩╚═╝ ",
	}

	c.BoxTopL, c.BoxTopC, c.BoxTopR = boxart.Row("╭──╮") // Top     ┌──┐ ╭──╮ ╒══╕ ╒══╕ ╔══╗
	c.BoxMidL, c.BoxMidC, c.BoxMidR = boxart.Row("│  │") // Mid     │  │ │  │ │  │ │  │ ║  ║
	c.BoxGndL, c.BoxGndC, c.BoxGndR = boxart.Row("├──┤") // Ground  ├──┤ ├──┤ ├──┤ ╞══╡ ╠══╣
	c.BoxBotL, c.BoxBotC, c.BoxBotR = boxart.Row("╰──╯") // Bottom  └──┘ ╰──╯ ╘══╛ ╘══╛ ╚══╝

	c.BoxInfoTL, c.BoxInfoT, c.BoxInfoTR = boxart.Row(" _ ")
	c.BoxInfoML, c.BoxInfoC, c.BoxInfoMR = boxart.Row("   ")
	c.BoxInfoBL, c.BoxInfoB, c.BoxInfoBR = boxart.Row(" ‾ ")

	c.BoxMidC = canvas

	c.TextAbc = "0123456789" + "`" +
		` -+*=/\.,:;!?$%&@#'"<>()[]{}^~_|` +
		`ABCDEFGHIJKLMNOPQRSTUVWXYZ` +
		`abcdefghijklmnopqrstuvwxyz`

	c.BlockAbc = strings.Split(c.TextAbc, "")
	return c
}

// TextToBlock overrides the default boxart.Boxboxart.TextToBlock lookup function
// with a direct passthrough function that returns the given str.
func (c *doublewidth) TextToBlock(str string) string {
	return str
}
