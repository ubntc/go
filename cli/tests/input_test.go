package tests

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/ubntc/go/cli/cli"
)

func TestInputChan(t *testing.T) {
	cli.GetTerm().SetDebug(true)

	f, remove := tempFile(t, "ab")
	defer remove()

	ch := cli.InputChan(f)
	i := 0
	for v := range ch {
		i++
		cli.Prompt("got rune %q", v)
	}

	assert.Equal(t, 2, i)
}

func TestProcessInput(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*500)

	cmds := []cli.Command{
		{Name: "command", Key: 'c', Fn: func() { cancel() }},
	}

	f, remove := tempFile(t, "cx")
	defer remove()

	go cli.ProcessInput(ctx, f, cmds)

	<-ctx.Done()
	assert.Equal(t, context.Canceled, ctx.Err())
}
