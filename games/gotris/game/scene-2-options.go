package game

import (
	"strconv"

	"github.com/ubntc/go/games/gotris/game/scenes"
)

func (g *Game) showOptions() {
	opt := &scenes.Scene{
		Name: scenes.Options,
	}
	for {
		names, idx := g.platform.RenderingModes()
		var infos []string
		for _, name := range names {
			infos = append(infos, g.platform.RenderingInfo(name))
		}
		opt.Options = names
		opt.Descriptions = infos
		opt.CurrentOption = idx

		key := g.ShowScene(opt, 0)

		switch cmd, _ := KeyToMenuCmd(key); cmd {
		case CmdEmpty, CmdMenuSelect:
			return
		case CmdMenuDown, CmdMenuRight:
			idx = (idx + 1) % len(names)
		case CmdMenuUp, CmdMenuLeft:
			idx = (idx + len(names) - 1) % len(names)
		}
		g.platform.SetRenderingMode(names[idx])
	}
}

func (g *Game) setRenderingMode(arg string) {
	names, _ := g.platform.RenderingModes()
	idx, _ := strconv.Atoi(arg)
	idx -= 1
	if idx >= len(names) {
		// ignore too high mode change requests
		return
	}
	g.platform.SetRenderingMode(names[idx])
}
