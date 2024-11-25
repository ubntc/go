package cli

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"

	xterm "golang.org/x/term"
)

var terminalLock sync.RWMutex

var (
	origStdout   = os.Stdout
	origStderr   = os.Stderr
	crPipeOut, _ = CrPipe(os.Stdout)
	crPipeErr, _ = CrPipe(os.Stderr)
)

func termMakeRaw() (*xterm.State, error) {
	terminalLock.Lock()
	defer terminalLock.Unlock()
	os.Stdout = crPipeOut
	os.Stderr = crPipeErr
	return xterm.MakeRaw(0)
}

func termRestore(state *xterm.State) error {
	terminalLock.Lock()
	defer terminalLock.Unlock()
	os.Stdout = origStdout
	os.Stderr = origStderr
	return xterm.Restore(0, state)
}

// RestoreFunc restores the terminal.
type RestoreFunc func() error

func restoreFunc(state *xterm.State) RestoreFunc {
	return func() error {
		if err := termRestore(state); err != nil {
			log.Println("failed to restore terminal, error:", err)
			return err
		}
		GetTerm().setRaw(false)
		log.Println("restored terminal")
		return nil
	}
}

// ClaimTerminal sets the terminal to raw input mode and returns a RestoreFunc
// for resetting the terminal to normal mode.
func ClaimTerminal() (RestoreFunc, error) {
	state, err := termMakeRaw()
	if err == nil {
		t := GetTerm()
		t.setRaw(true)
		log.Println("claimed terminal")
		if t.IsDebug() && t.IsVerbose() {
			t.Println(fmt.Sprintf("---termios state---\n\r%+v\n\r------", state))
		}
	} else {
		log.Println("failed to claim terminal")
	}

	var restore RestoreFunc
	// If terminal is in raw mode now, let's ensure it can return to normal.
	if state != nil {
		restore = restoreFunc(state)
		if err != nil {
			log.Println("Warning: terminal RestoreFunc will be using a nil xterm.State")
		}
	}

	return restore, err
}

// InputChan listens for runes on stdin
// and writes them to the returned channel.
func InputChan(file *os.File) <-chan rune {
	ch := make(chan rune, 10)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				Prompt("Failed to read from '%s', panic=%v", file.Name(), r)
			}
		}()
		in := bufio.NewReader(file)
		defer close(ch)
		for {
			// this call can panic if the stdin pipeline is augmented
			char, _, err := in.ReadRune()
			if err != nil && err != io.EOF {
				Prompt("Reader failed, error=%v", err)
				return
			}
			if err != nil {
				debug("reader stopped")
				return
			}
			debug("read char: %q", char)
			ch <- char
		}
	}()
	return ch
}

// ProcessInput reads runes from input chan and executes the `commands` mapped to the received input keys.
func ProcessInput(ctx context.Context, file *os.File, commands Commands, termMakeRaw bool) {
	debug("start processing input")
	defer debug("input processing stopped")

	// when reading from stdin, acquire raw terminal input and make ProcessInput wait for terminal after cleanup
	if file == os.Stdin && termMakeRaw {
		restore, err := ClaimTerminal()
		if restore != nil {
			defer restore() // nolint
		}
		if err != nil {
			debug("failed to claim terminal, error=%s", err.Error())
		}
	}

	input := InputChan(file)

	Prompt(commands.String())

	var prompt string
	var char rune
	var more bool
	for {
		select {
		case <-ctx.Done():
			debug("Quit (context done).")
			return
		case <-time.After(time.Second):
			if len(prompt) > 0 {
				Prompt(prompt)
				prompt = ""
			}
		case char, more = <-input:
			if !more {
				<-ctx.Done()
				debug("Quit (input closed + context done).")
				return
			}
			if cmd := commands.Get(char); cmd != nil {
				prompt = ""
				go cmd.Run(ctx)
				continue
			}
			Prompt("Pressed key %q.", char)
			prompt = commands.String()
		}
	}
}
