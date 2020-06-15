package stdlogger

import (
	"fmt"
	"io"
	"log"
	"time"
)

const (
	colorYellow   = 33
	colorDarkGray = 90
)

// colored returns a colored string.
func colored(s interface{}, c int) string {
	return fmt.Sprintf("\x1b[%dm%v\x1b[0m", c, s)
}

// ConsoleWriter writes to the console.
type ConsoleWriter struct {
	out        io.Writer
	timeFormat string
}

// Write mimics zerolog.ConsolWrite Output.
func (c *ConsoleWriter) Write(p []byte) (n int, err error) {
	t := colored(time.Now().Format(c.timeFormat), colorDarkGray)
	s := colored("DBG", colorYellow)
	return c.out.Write([]byte(fmt.Sprintf("%s %s %s", t, s, p)))
}

// Setup sets up the standard logger.
func Setup(out io.Writer, timeFormat string) error {
	// setup a ConsoleWriter to add colored timestamp and level
	log.SetOutput(&ConsoleWriter{out, timeFormat})
	// remove std. logger time flags, use only prefix flag
	log.SetFlags(log.Lmsgprefix)
	log.Println("using stdlogger.ConsoleWriter")
	return nil
}
