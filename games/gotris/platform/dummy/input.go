package fyne

import (
	"context"
	"os"

	"github.com/ubntc/go/games/gotris/input"
	"github.com/ubntc/go/games/gotris/terminal"
)

func (p *Platform) CaptureInput(ctx context.Context) (<-chan input.Key, func(), error) {
	t := terminal.New(os.Stdout)
	return t.CaptureInput(ctx)
}
