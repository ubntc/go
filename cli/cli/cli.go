package cli

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/ubntc/go/cli/cli/config"
)

func newSigChan() <-chan os.Signal {
	// size > 1 makes the channel non-blocking
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	return sig
}

// SigWait waits for OS signals SIGINT or SIGTERM or the termination of the given context.
// It blocks until the context is canceled either by the awaited signal or externally.
// It returns the received signal and the context's error.
//
// Usage:
//
//		 ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
//		 defer cancel()
//	     // start async workloads
//	     go myServer.Start(ctx)
//	     // await programm termination
//		 sig, err := cli.SigWait(ctx)
//		 fmt.Println("stopping application")
func SigWait(ctx context.Context) (os.Signal, error) {
	sig := newSigChan()
	var s os.Signal

	// block until signal is received or context is cancelled
	select {
	case <-ctx.Done():
	case s = <-sig:
	}

	// check and return context errors (excl. cancellation)
	// this way the caller can check for abnornal program behavior
	if err := ctx.Err(); err != context.Canceled {
		return s, err
	}
	return s, nil
}

// StartTerm starts terminal session management and input handling (if configured)
// and returns a context.Context to manage the corresponding goroutines or running i/o pipes.
// It cancels the context on receiving a SIGINT or SIGTERM from the OS.
func StartTerm(parent context.Context, cfg config.Config, cmds ...Command) (context.Context, context.CancelFunc) {
	sig := newSigChan()
	input := os.Stdin
	// var interactive bool

	// use separate contexts to wait for closing input and for the clock
	// use separate cancel functions to trigger stop reading input and clock stop
	inputCtx, stopReadingInput := context.WithCancel(context.Background())
	clockCtx, stopClock := context.WithCancel(context.Background())

	// create a new context to be returned that is disconnected from parent
	// to allow for sigwait to finish processing and cleanup the terminal
	ctx, cancel := context.WithCancel(context.Background())

	opts := cfg
	var commands Commands

	commands = append(commands, cmds...)

	if opts.WithQuit {
		commands = append(commands, QuitCommands(stopReadingInput)...)
	}

	// log.Println("added", len(commands), "commands")

	// wg blocks canceling the returned context until after:
	// 1. input processing has stopped
	// 2. the clock has stopped
	var wg sync.WaitGroup

	// start clock separately
	if opts.ShowClock {
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer stopClock()
			GetTerm().StartClock(clockCtx) // blocking
		}()
	}

	restoreStdio := func() {}
	if cfg.PrependCR {
		os.Stderr = crPipeErr
		os.Stdout = crPipeOut
		restoreStdio = func() {
			os.Stderr = origStderr
			os.Stdout = origStdout
		}
	}

	// start reading input separately
	// manage the terminal only if there are some commands to handle
	if len(commands) > 0 {
		// update global commands
		SetCommands(commands)

		if cfg.ShowClock && !cfg.MakeTermRaw {
			log.Println("⚠️ setting MakeTermRaw=true ⚠️ (required by ShowClock in interative mode)")
			cfg.MakeTermRaw = true
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer stopReadingInput()
			ProcessInput(inputCtx, input, GetCommands(), cfg.MakeTermRaw) // blocking
		}()
	}

	go func() {
		defer func() {
			cancel() // cancel the exposed context
			debug("main context cancelled")

			stopReadingInput() // ensure inputCtx is also canceled
			stopClock()        // ensure clockCtx is also canceled
			restoreStdio()     // ensure os.Stderr and os.Stderr are restored

			debug("wait for cleanup")
			wg.Wait()
			debug("cleanup finished")
		}()
		select {
		case s := <-sig:
			debug("stop on signal: %q", s)
		case <-parent.Done():
			debug("stop on closing parent")
		case <-inputCtx.Done():
			debug("stop on closing input")
		case <-clockCtx.Done():
			debug("stop on stopped clock")
		}
	}()

	return ctx, cancel
}
