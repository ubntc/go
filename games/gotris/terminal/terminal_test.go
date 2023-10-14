package terminal

import (
	"os"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClear(t *testing.T) {
	assert.Contains(t, clear, runtime.GOOS)
	// cannot call clear functions in test, because it would clear test output
	// TODO: use dummy file descriptor

	dummyout, err := os.Open(os.DevNull)
	assert.NoError(t, err)
	term := NewTerminal(dummyout)

	term.Print(term.ClearString())
	term.RunClearCommand()

	/*
		// TODO: How to test?
		// Error: inappropriate ioctl for device, cannot overpaint, failed to get size
		err = term.Overpaint()
		assert.NoError(t, err)

		// TODO: How to test CaptureInput?
		// * Cannot capture input during `go test` run.
		// * Cannot use dummy input.
		dummyin, err := os.Open(os.DevNull)
		assert.NoError(t, err)
		term.stdin = dummyin

		ctx, cancel := context.WithCancel(context.Background())
		ch, restore, err := term.CaptureInput(ctx)
		defer restore()
		assert.NoError(t, err)
		cancel()
		<-ch
	*/
}

func TestTerminal(_ *testing.T) {
	t := NewTerminal(os.Stdout)
	defer t.ShowCursor()
	t.HideCursor()
	t.Print("test")
	t.Println("test")
}
