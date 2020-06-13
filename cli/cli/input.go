package cli

import (
	"bufio"
	"context"
	"io"
	"os"
	"sync"
	"time"

	"golang.org/x/crypto/ssh/terminal"
)

var terminalLock sync.RWMutex

func termMakeRaw() (*terminal.State, error) {
	terminalLock.Lock()
	defer terminalLock.Unlock()
	return terminal.MakeRaw(0)
}

func termRestore(state *terminal.State) error {
	terminalLock.Lock()
	defer terminalLock.Unlock()
	return terminal.Restore(0, state)
}

// RestoreFunc restores the terminal.
type RestoreFunc func() error

// ClaimTerminal sets the terminal to raw input mode, waits for the given
// context to cancel, and eventually resets the terminal to normal mode.
func ClaimTerminal() (RestoreFunc, error) {
	PromptVerbose("using raw terminal")
	state, err := termMakeRaw()
	if err != nil {
		return nil, err
	}
	GetTerm().setRaw(true)

	// Terminal is in raw mode now. Let's ensure it can return to normal.
	return func() error {
		PromptVerbose("restoring terminal")
		if err := termRestore(state); err != nil {
			return err
		}
		GetTerm().setRaw(false)

		return nil
	}, nil
}

// InputChan listens for runes on stdin
// and writes them to the returned channel.
func InputChan(file *os.File) <-chan rune {
	ch := make(chan rune, 10)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				Prompt("Failed to read from %s, panic=%v", file.Name(), r)
			}
		}()
		in := bufio.NewReader(file)
		defer close(ch)
		for {
			// this call can panic if the stdin pipeline is augmented
			char, _, err := in.ReadRune()
			if err != nil && err != io.EOF {
				Prompt("Reader failed, error=%v", err)
			}
			if err != nil {
				PromptVerbose("reader stopped")
				return
			}
			PromptVerbose("read char: %q", char)
			ch <- char
		}
	}()
	return ch
}

// ProcessInput reads runes from input chan and executes the `commands` mapped to the received input keys.
func ProcessInput(ctx context.Context, file *os.File, commands Commands) {
	// when reading from stdin, acquire raw terminal input and make ProcessInput wait for terminal after cleanup
	if file == os.Stdin {
		restore, err := ClaimTerminal()
		if err != nil {
			PromptVerbose("failed to switch to terminal mode, error=%s", err)
		} else {
			defer restore()
		}
	}

	var wg sync.WaitGroup
	defer wg.Wait() // block ProcessInput to ensure terminal cleanup

	input := InputChan(file)

	Prompt(commands.String())

	var prompt string
	var char rune
	var more bool
	for {
		select {
		case <-ctx.Done():
			PromptVerbose("Quit (context done).")
			return
		case <-time.After(time.Second):
			if len(prompt) > 0 {
				Prompt(prompt)
				prompt = ""
			}
		case char, more = <-input:
			if !more {
				<-ctx.Done()
				PromptVerbose("Quit (input closed + context done).")
				return
			}
			if cmd := commands.Get(char); cmd != nil {
				prompt = ""
				wg.Add(1)
				go func() {
					defer wg.Done()
					cmd.Run()
				}()
				continue
			}
			Prompt("Pressed key %q.", char)
			prompt = commands.String()
		}
	}
}
