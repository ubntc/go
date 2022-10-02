package terminal

import (
	"fmt"
	"os"
	"strings"

	"github.com/pkg/errors"
	xterm "golang.org/x/term"
)

const (
	EscClearScreen = "\x1b[2J"
	EscGotoTopLeft = "\x1b[0;0f"
	EscHideCursor  = "\x1b[?25l"
	EscShowCursor  = "\x1b[?25h"
)

type Terminal struct {
	stdout *os.File
	stdin  *os.File
}

// New returns a new terminal for thr given file descriptor.
func New(stdout *os.File) *Terminal {
	return &Terminal{stdout, os.Stdin}
}

// Overpaint clears the screen by printing w*h spaces.
func (t *Terminal) Overpaint() error {
	w, h, err := xterm.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return errors.Wrap(err, "cannot overpaint, failed to get size")
	}
	line := strings.Repeat(" ", w)

	t.Print(EscGotoTopLeft)
	for x := 0; x < h-1; x++ {
		t.Println(line)
	}
	if os.Getenv("DEBUG") != "" {
		info := fmt.Sprintf("w:%d, h:%d", w, h)
		info += strings.Repeat(" ", w-len(info))
		t.Println(info)
	}
	t.Print(EscGotoTopLeft)
	return nil
}

// RunClearCommand runs a "clear" or similar command if supported by the OS.
// Panics if the OS does not support "clear".
func (t *Terminal) RunClearCommand() {
	// works in iTerm2 with some flickering and perfectly in VSCode terminal
	callClearFunc(t.stdout)
}

// Clear sends ANSI codes for clearing the screen and setting the cursor to 0,0.
func (t *Terminal) Clear() {
	// works in iTerm2 and VSCode terminal with minor flickering
	t.Print(EscGotoTopLeft + EscClearScreen + EscGotoTopLeft)
}

// Print prints the values to the terminal's file descriptor.
func (t *Terminal) Print(values ...interface{}) {
	fmt.Fprint(t.stdout, values...)
}

// Println prints the values to the terminal's file descriptor ending with a new line
// and an additional carriage return "\r", to produce valid lines on raw terminals.
func (t *Terminal) Println(values ...interface{}) {
	fmt.Fprint(t.stdout, values...)
	fmt.Fprintln(t.stdout)
	fmt.Fprint(t.stdout, "\r") // needed for raw terminals
}

// ShowCursor sends ANSI escape code to show the cursor.
func (t *Terminal) ShowCursor() {
	fmt.Fprint(t.stdout, EscShowCursor)
}

// HideCursor sends ANSI escape code to hide the cursor.
func (t *Terminal) HideCursor() {
	fmt.Fprint(t.stdout, EscHideCursor)
}
