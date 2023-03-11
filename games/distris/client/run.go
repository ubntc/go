package client

import (
	"context"
	"fmt"
	"os"

	"github.com/ubntc/go/games/distris/api/command"
	"github.com/ubntc/go/games/gotris/terminal"
)

func Run(address string) error {
	c := New(address)
	fmt.Printf("connecting to server: %s\n\r", address)
	term := terminal.New(os.Stdin)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	keys, restore, err := term.CaptureInput(ctx)
	if restore != nil {
		defer restore()
	}
	if err != nil {
		return err
	}

	for {
		select {
		case key, more := <-keys:
			if key != nil {
				if err := c.Send(ctx, command.Command(key.Text())); err != nil {
					return err
				}
			}
			if !more {
				cancel()
			}
		case <-ctx.Done():
			fmt.Println("client stopped")
			return nil
		}
	}
}
