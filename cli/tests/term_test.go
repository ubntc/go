package tests

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/ubntc/go/cli/cli"
)

func TestPrompt(t *testing.T) {
	term := cli.GetTerm()
	term.SetVerbose(true)
	defer term.SetVerbose(false)

	cli.Prompt("test")
	assert.Equal(t, "test", term.GetMessage())

	cli.PromptVerbose("test verbose")
	assert.Equal(t, "test verbose", term.GetMessage())
}

func TestPrint(t *testing.T) {
	term := cli.GetTerm()
	i, err := term.Println("test")
	assert.NoError(t, err)
	assert.Equal(t, 5, i)
}

func TestHelp(t *testing.T) {
	cli.SetCommands(cli.Commands{cli.Command{Name: "X", Key: 'x'}})
	term := cli.GetTerm()
	w := &LogWriter{}
	term.WrapOutput(w)
	term.Help()
	s := w.String()
	assert.Contains(t, s, "Key: 'x'")
	assert.Contains(t, s, "Command: X")
}

func TestQuitCommandsWithoutClock(t *testing.T) {
	ctx, cancel := cli.WithSigWait(context.Background(), cli.WithQuit(), cli.WithoutClock())
	assert.Len(t, cli.GetCommands(), 4)
	cancel()
	<-ctx.Done()
}

func TestCommandsWithClock(t *testing.T) {
	pctx, pcancel := context.WithTimeout(context.Background(), time.Millisecond*100)
	defer pcancel()

	cmds := cli.Commands{
		{Name: "cancel", Key: 'c', Fn: func() { pcancel() }},
	}

	term, release := cli.AcquireTerm()
	defer release()
	w := &LogWriter{}
	term.WrapOutput(w)

	ctx, cancel := cli.WithSigWait(pctx, cli.WithInput(cmds))
	defer cancel()

	assert.Len(t, cli.GetCommands(), 5)
	time.Sleep(20 * time.Millisecond)
	// run the command as soon as the clock gets active
	cmds.Run('c')
	assert.Contains(t, w.String(), time.Now().Format("2006-01-02"))

	<-ctx.Done()
	assert.Equal(t, context.Canceled, ctx.Err())
}

func TestWrite(t *testing.T) {
	assert.Equal(t, string(cli.CR), "\r")
	assert.Equal(t, string(cli.NL), "\n")

	term, release := cli.AcquireTerm()
	defer release()
	i, err := term.Write([]byte("TestWrite"))
	assert.NoError(t, err, "term.Write should not produce an error")
	assert.Equal(t, 9, i, "9 bytes should be written")
	i, err = term.Sync()
	assert.NoError(t, err, "term.Write should not produce an error")
	assert.Equal(t, 1, i, "1 sync byte should be written")
	term.Sync()
}

var reLastNonEmptyLine = regexp.MustCompile("[^\n\r]*[\n][\r]*$")

// test-write message to LogWriter and test if it was wrtten correctly.
func write(t *testing.T, w *LogWriter, s string, output string) {
	term := cli.GetTerm()
	i, err := term.WriteString(s)
	// println("buffer:", w.Quote())
	// println("raw:", term.IsRaw())
	lastLine := reLastNonEmptyLine.FindString(w.String())
	if len(lastLine) > 0 {
		lastIdx := len(lastLine) - 1
		if lastLine[lastIdx] == cli.CR {
			lastLine = lastLine[0:lastIdx]
		}
	}
	assert.Equal(t, false, cli.GetTerm().IsRaw())
	assert.NoError(t, err, "term.Write should not produce an error")
	assert.Equal(t, len(s), i, "all bytes should be written")
	assert.Equal(t, output, lastLine, "output should match", w.Quote(), term.IsRaw())
}

func TestWrappedOutput(t *testing.T) {
	term, release := cli.AcquireTerm()
	defer release()
	w := &LogWriter{}
	term.WrapOutput(w)

	// fmt.Printf("terminal raw mode: %v\n", term.IsRaw())

	write(t, w, "test", "")
	write(t, w, "\n", "test\n")
	w.Flush()

	write(t, w, "", "")
	w.Flush()

	write(t, w, "123", "")
	write(t, w, "456\n", "123456\n")
	write(t, w, "789\n", "789\n")
	w.Flush()

	// TODO: add clock tests
	term.Sync()
}

func TestWrappedOutputConcurrent(t *testing.T) {
	term, release := cli.AcquireTerm()
	defer release()
	w := &LogWriter{}
	term.WrapOutput(w)
	var wg sync.WaitGroup

	flush := func() { term.Sync() }
	write := func(s string) { term.WriteString(s) }

	worker := func(fn func()) {
		defer wg.Done()
		stop := time.After(10 * time.Millisecond)
		for {
			select {
			case <-stop:
				return
			case <-time.After(time.Millisecond * 3):
				wg.Add(1)
				go func() {
					defer wg.Done()
					fn()
				}()
			}
		}
	}

	assert.Equal(t, false, cli.GetTerm().IsRaw())

	wg.Add(4)
	go worker(func() { write("a") })
	go worker(func() { write("b") })
	go worker(func() { write("c\n") })
	go worker(func() { write("d\n") })
	wg.Wait()

	wg.Add(2)
	go worker(flush)
	go worker(flush)
	wg.Wait()

	assert.GreaterOrEqual(t, len(w.String()), 8)

	res := w.String()
	assert.Regexp(t, "^(a|b|c\\n|d\\n|\\r|\\n| )*$", res)

	term.Sync()
}

func TestWrapStderr(t *testing.T) {
	v := Capture(os.Stdout, func() {
		fmt.Println("test")
	})
	assert.Equal(t, v, "test\n")
}

func TestSliceIndexPlus1(t *testing.T) {
	var buf = []byte("abc")
	assert.Len(t, buf, 3)
	assert.Len(t, buf[len(buf):], 0)
}
