package cli

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type options struct {
	withoutClock bool
	withQuit     bool
	commands     []Command
}

// Option defines a sigwait option.
type Option func(*options)

// WithInput makes sigwait run the given commands on receiving user input.
// The keys q, Q, CTRL-C, and CTRL-D are reserved to quit the program.
func WithInput(commands []Command) Option {
	return func(o *options) {
		o.commands = append(o.commands, commands...)
	}
}

// WithQuit add the default quit commands and enables user input.
func WithQuit() Option {
	return func(o *options) {
		o.withQuit = true
	}
}

// WithoutClock disabled the default ascii/unicode clock in the last terminal line.
func WithoutClock() Option {
	return func(o *options) {
		o.withoutClock = true
	}
}

// SigWait waits for OS signals and cancels the given context on SIGINT or SIGTERM.
// It blocks until the context is canceled either by the awaited signal or externally.
// It returns the received signal and the context's error.
func SigWait(ctx context.Context, cancel context.CancelFunc) (os.Signal, error) {
	defer cancel()
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	var s os.Signal

	select {
	case <-ctx.Done():
	case s = <-sig:
	}

	if err := ctx.Err(); err != context.Canceled {
		return s, err
	}
	return s, nil
}

// WithSigWait returns a context.Context that is canceled on OS signal.
//
// WithSigWait starts a goroutine that waits for OS signals SIGINT or SIGTERM and cancels the
// returned context after receiving the signal. Depending on the options, WithSigWait starts
// processing input and sets up commands.
//
func WithSigWait(parent context.Context, opt ...Option) (context.Context, context.CancelFunc) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	input := os.Stdin
	var inputCommands Commands
	var interactive bool

	// use separate contexts to wait for closing input and for the clock
	// use separate cancel functions to trigger stop reading input and clock stop
	inputCtx, stopReadingInput := context.WithCancel(context.Background())
	clockCtx, stopClock := context.WithCancel(context.Background())

	// create a new context to be returned that is disconnected from parent
	// to allow for sigwait to finish processing and cleanup the terminal
	sigWaitCtx, gracefulCancel := context.WithCancel(context.Background())

	var opts options
	for _, apply := range opt {
		apply(&opts)
	}

	if len(opts.commands) > 0 || opts.withQuit {
		inputCommands = QuitCommands(stopReadingInput)
		inputCommands = append(inputCommands, opts.commands...)
		interactive = true
	}

	// wg blocks canceling the returned context until after:
	// 1. input processing has stopped
	// 2. the clock has stopped
	var wg sync.WaitGroup

	if interactive && !opts.withoutClock {
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer stopClock()
			PromptVerbose("clock started")
			GetTerm().StartClock(clockCtx)
			PromptVerbose("clock finished")
		}()
	}

	if interactive {
		// update global commands
		SetCommands(inputCommands)
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer stopReadingInput()
			PromptVerbose("user input started")
			ProcessInput(inputCtx, input, GetCommands())
			PromptVerbose("user input finished")
		}()
	}

	go func() {
		defer gracefulCancel() // cancel returned context after all cleanup is done
		select {
		case s := <-sig:
			PromptVerbose("stop on signal: %q", s)
		case <-parent.Done():
			PromptVerbose("stop on closing parent")
		case <-inputCtx.Done():
			PromptVerbose("stop on closing input")
		case <-clockCtx.Done():
			PromptVerbose("stop on stopped clock")
		}
		stopReadingInput() // ensure inputCtx is canceled
		stopClock()        // ensure clockCtx is canceled
		PromptVerbose("wait for cleanup")
		wg.Wait() // wait for cleanup
		PromptVerbose("cleanup finished")
	}()

	return sigWaitCtx, gracefulCancel
}
