package cli

import (
	"context"
	"strings"
	"sync"
	"time"
)

// Clock stores the current clock step and text.
type Clock struct {
	timeFormat string
	runLock    sync.Mutex // ensures only one clock runs concurrently
}

// DisplayTime defines the visible values for a time.
type DisplayTime struct {
	time.Time
	digital string
	analog  string
}

var (
	// asciiClock     = strings.Split("/ - \\ |", " ")
	brailleClock = strings.Split("⢎⡰ ⢎⡡ ⢎⡑ ⢎⠱ ⠎⡱ ⢊⡱ ⢌⡱ ⢆⡱", " ")
	// brailleSpinner = strings.Split(" ⠁| ⠑| ⠰| ⡰|⢀⡠|⢄⡠|⢆⡀|⢎⡀|⢎ |⠎ |⠊ |⠈ ", "|")
)

// Start starts the clock.
func (c *Clock) Start(ctx context.Context, tickTime time.Duration) <-chan *DisplayTime {
	c.runLock.Lock()
	out := make(chan *DisplayTime)

	go func() {
		defer c.runLock.Unlock()
		defer close(out)
		for {
			out <- c.DisplayTime(tickTime)
			select {
			case <-time.After(tickTime):
			case <-ctx.Done():
				return
			}
		}
	}()
	return out
}

// Chars returns the clock chars.
func (c *Clock) Chars(tickInterval time.Duration) string {
	step := (time.Now().UnixNano() / int64(tickInterval)) % int64(len(brailleClock))
	return brailleClock[step%int64(len(brailleClock))]
}

// DisplayTime returns the current clock text.
func (c *Clock) DisplayTime(tickInterval time.Duration) *DisplayTime {
	t := time.Now()
	return &DisplayTime{t, t.Format(c.timeFormat), c.Chars(tickInterval)}
}
