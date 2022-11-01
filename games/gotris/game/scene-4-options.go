package game

import (
	cmd "github.com/ubntc/go/games/gotris/game/controls"
	"github.com/ubntc/go/games/gotris/game/scenes"
)

func (g *Game) showOptions() {
	scn := scenes.NewOptionsScene(scenes.TitleOptions, g.platform.Options())
	opt := scn.Options()

	if len(opt.List()) == 0 {
		return
	}

	for {
		key := g.ShowScene(scn, 0)
		c, _ := cmd.KeyToMenuCmd(key)
		if i, done, ok := cmd.HandleOptionsCmd(c, len(opt.List()), opt.Get()); ok {
			opt.Set(i)
			if done {
				break
			}
		}
	}
}
