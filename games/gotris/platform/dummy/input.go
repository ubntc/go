package dummy

import (
	"context"
	"os"

	"github.com/ubntc/go/games/gotris/input"
	"github.com/ubntc/go/games/gotris/terminal"
)

func (p *Platform) CaptureInput(ctx context.Context) (<-chan input.Input, func(), error) {
	t := terminal.New(os.Stdout)
	return t.CaptureInput(ctx)
}
