package game

import (
	"github.com/ubntc/go/games/gotris/common/labels"
	"github.com/ubntc/go/games/gotris/common/platform"
	cmd "github.com/ubntc/go/games/gotris/game/controls"
)

func (g *Game) showOptions() {
	scn := platform.NewMenu(labels.TitleOptions, g.Platform.Options())
	opt := scn.Options()

	if len(opt.List()) == 0 {
		return
	}

	for {
		key := g.ShowScene(scn, 0)
		c, _ := cmd.KeyToMenuCmd(key)
		if result := cmd.HandleOptionsCmd(c, opt); result == cmd.HandleResultSelectionFinished {
			return
		}
	}
}
