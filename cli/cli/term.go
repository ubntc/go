package cli

import (
	"context"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	xterm "golang.org/x/term"
)

// Term stores shared terminal state for loggers, the underlying terminals, and commands.
type Term struct {
	out        io.Writer
	debug      bool
	verbose    bool
	raw        bool
	mu         sync.RWMutex
	clock      clock
	commands   Commands
	buf        []byte // last received bytes received from external writers
	statusLine string // current status line text
	message    string // message to be displayed on the status line
	lastLine   string // last line that was printed
}

// the global term
var term = Term{
	out: os.Stderr,
	// the global Clock
	clock: Clock(clockClock),
	// global commands
	commands: Commands{},
}

// lock to acquire global term
var aqmu sync.Mutex

// GetTerm returns the global terminal state.
func GetTerm() *Term {
	return &term
}

// AcquireTerm locks and returns the global terminal.
func AcquireTerm() (*Term, func()) {
	aqmu.Lock()
	return GetTerm(), aqmu.Unlock
}

// IsDebug returns the debug state.
func (c *Term) IsDebug() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.debug
}

// SetDebug enabled or disables debug output on stderr.
func (c *Term) SetDebug(v bool) *Term {
	c.mu.Lock()
	defer c.mu.Unlock()
	term.debug = v
	return c
}

// IsVerbose returns the verbose state.
func (c *Term) IsVerbose() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.verbose
}

// SetVerbose enabled or disables verbose output on stderr.
func (c *Term) SetVerbose(v bool) *Term {
	c.mu.Lock()
	defer c.mu.Unlock()
	term.verbose = v
	return c
}

// SetMessage set the promt message.
func (c *Term) SetMessage(msg string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	term.message = msg
}

// GetMessage set the promt message.
func (c *Term) GetMessage() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return term.message
}

// Help prints and prompts help for the configured Commands.
func (c *Term) Help() {
	c.SetMessage(GetCommands().String())
	c.Println(GetCommands().Help())
}

// SetOutput set the Term's output.
func (c *Term) SetOutput(w io.Writer) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.out = w
}

// var reLineEnd = regexp.MustCompile("\n$")
var (
	rePendingNL   = regexp.MustCompile("[^\n]*$")
	reStartCR     = regexp.MustCompile("^[\r]*")
	reNLCR        = regexp.MustCompile("[\n\r]*")
	reNLWithoutCR = regexp.MustCompile("[\n][^\r]*")
)

// TODO: what is safer/faster regex check or last bytes check?
// var nlByte = []byte("\n")[0]
// var crByte = []byte("\r")[0]

func (c *Term) Write(b []byte) (int, error) {
	if len(b) == 0 {
		return 0, nil
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.buf = append(c.buf, b...)
	c.write()
	return len(b), nil
}

// WriteString writes the string to the terminal.
func (c *Term) WriteString(s string) (int, error) {
	return c.Write([]byte(s))
}

// Println writes the string + "\n" to the terminal.
func (c *Term) Println(s string) (int, error) {
	return c.Write([]byte(s + "\n"))
}

// printableOutput returns completed and pending lines in the output buffer.
func (c *Term) printableOutput() (output, pending []byte) {
	locNL := rePendingNL.FindIndex(c.buf)
	if locNL == nil {
		return c.buf, nil
	}
	i := locNL[0]
	locCR := reStartCR.FindIndex(c.buf[i:])
	if locCR == nil {
		i = locNL[0] + locCR[0]
	}
	return c.buf[0:i], c.buf[i:]
}

// Sync flushes buffers (appending newlines if needed) and clears all output.
func (c *Term) Sync() (int, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if len(c.buf) == 0 {
		return 0, nil
	}
	c.write()
	if len(c.buf) > 0 {
		c.buf = append(c.buf, []byte("\n")...)
		c.write()
		return 1, nil
	}
	c.buf = nil
	return 0, nil
}

// Prompt prints the prompt string in the termnial line.
func (c *Term) Prompt(v ...string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	s := strings.Join(v, " ")
	s = reNLCR.ReplaceAllLiteralString(s, "")
	term.statusLine = s
	c.write()
}

const (
	// CR is the carriage return byte
	CR = byte(13)

	// NL is the newline byte
	NL = byte(10)
)

// The internal `write` writes the collected `input` followed by the `statusLine`.
// The following cases need to be handled:
//
//	input buffer    prompt status    action
//	---------------------------------------
//	empty           unchanged        non
//	empty           updated          clear + print status
//	has data        any              clear + print input and status
func (c *Term) write() {
	output, pending := c.printableOutput()
	if len(output) == 0 && c.statusLine == c.lastLine {
		return
	}
	var buf []byte
	buf = append(buf, []byte(c.clearString())...)
	buf = append(buf, output...) // output is nil or ends with NL/CR

	// CR check and handling in raw mode
	end := buf[len(buf)-1]
	if c.debug && end != CR && end != NL {
		panic(fmt.Errorf("invalid buffer data: %q", buf))
	}
	if c.raw {
		// TODO: do we need to add CR for each inline NL?
		// Yes! ;)
		locs := reNLWithoutCR.FindAllIndex(buf, -1)
		if len(locs) > 0 {
			b := make([]byte, len(buf)+len(locs))
			loc1 := 0
			for _, loc := range locs {
				b = append(b, buf[loc1:loc[0]]...)
				b = append(b, []byte("\n\r")...)
				loc1 = loc[0] + 1
				// Since loc[0] is always a valid byte, loc[0] + 1 is at most == len(buf),
				// which is a valid range index for a slice
			}
			buf = append(b, buf[loc1:]...)
			if len(locs) > 1 && c.debug && c.verbose {
				// Note: len(locs) == 1 is the default case, no need to log it.
				buf = append(buf, []byte(fmt.Sprintf("injected %d CR\n\r", len(locs)))...)
			}
		}
	}

	buf = append(buf, []byte(c.statusLine)...)
	_, _ = c.out.Write(buf)
	c.buf = pending
	c.lastLine = c.statusLine
}

// clearString returns a string to clear the complete line.
func (c *Term) clearString() string {
	if w, _, err := xterm.GetSize(0); err == nil && w > 0 {
		return fmt.Sprintf("\r%s\r", strings.Repeat(" ", w))
	}
	return ClearAll
}

// SetRaw sets the raw state.
func (c *Term) setRaw(v bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	term.raw = v
}

// IsRaw returns the raw state.
func (c *Term) IsRaw() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return term.raw
}

// StartClock starts the terminal clock to ingest clock output into the status line.
func (c *Term) StartClock(ctx context.Context) {
	debug("starting clock")
	interval := 100 * time.Millisecond
	ticker := time.NewTicker(interval)
	for {
		select {
		case <-ticker.C:
			dt := c.clock.DisplayTime(interval)
			if dt != nil {
				c.Prompt(dt.digital, c.GetMessage(), dt.analog, "")
			}
		case <-ctx.Done():
			debug("clock stopped")
			c.Prompt("")
			return
		}
	}
}

// GetClock returns the global clock.
func (c *Term) GetClock() *clock {
	return &c.clock
}
