package tests

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/ubntc/go/cli/cli"
)

// LogWriter implements io.Writer.
type LogWriter struct {
	sync.RWMutex
	value []byte
}

func (w *LogWriter) Write(b []byte) (int, error) {
	w.Lock()
	defer w.Unlock()
	w.value = append(w.value, b...)
	return len(b), nil
}

// Flush flushes the write buffer
func (w *LogWriter) Flush() {
	w.Lock()
	defer w.Unlock()
	w.value = nil
}

func (w *LogWriter) String() string {
	w.RLock()
	defer w.RUnlock()
	return string(w.value)
}

// Quote returns the quoted output.
func (w *LogWriter) Quote() string {
	w.RLock()
	defer w.RUnlock()
	return fmt.Sprintf("%q", w.value)
}

var tfLen = len(cli.TimeFormatHuman)

// var clrMatch = regexp.MustCompile("^\r *\r$")
var emptyMatch = regexp.MustCompile("^ *$")

// Time returns the time in the current output line.
func (w *LogWriter) Time() (time.Time, error) {
	text := w.String()
	var s, line string
lines:
	for _, line = range strings.Split(text, "\n") {
		// skip empty lines
		if emptyMatch.Match([]byte(line)) {
			continue
		}
		for _, s = range strings.Split(line, "\r") {
			// skip empty clear strings
			if emptyMatch.Match([]byte(s)) {
				continue
			}
			// stop as soon as we find a non empty string
			break lines
		}
	}
	if len(s) < tfLen {
		return time.Time{}, fmt.Errorf("time string to short: s='%s', line='%s'", s, line)
	}
	return time.Parse(cli.TimeFormatHuman, s[:tfLen])
}

// TempFile creates a temp file.
func TempFile(t *testing.T, content string) (*os.File, func()) {
	f, err := os.CreateTemp(os.TempDir(), "gocli")
	remove := func() {
		assert.NoError(t, os.Remove(f.Name()))
	}
	assert.NoError(t, err)
	_, err = f.WriteString(content)
	assert.NoError(t, err)
	_, err = f.Seek(0, 0)
	assert.NoError(t, err)
	return f, remove
}
