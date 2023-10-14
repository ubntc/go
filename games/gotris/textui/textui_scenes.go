package textui

import (
	"strings"

	"github.com/ubntc/go/games/gotris/common/labels"
	"github.com/ubntc/go/games/gotris/common/platform"
	txt "github.com/ubntc/go/games/gotris/textui/textscenes"
)

func (ui *TextUI) RenderScene(scene platform.Scene) {
	var screen string
	opt := scene.Options()
	switch scene.Name() {
	case labels.TitleWelcome:
		screen = txt.WelcomeScreen.Menu(opt.List(), opt.Descs(), opt.GetName())
	case labels.TitleOptions:
		screen = txt.OptionScreen.Menu(opt.List(), opt.Descs(), opt.GetName())
	case labels.TitleControls:
		screen = txt.Controls
	case labels.TitleGameOver:
		screen = txt.GameOver
	}
	ui.Clear()
	lines := strings.Split(string(screen), "\n")
	ui.Print(strings.Join(lines, "\r\n    "))
}
