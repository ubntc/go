package cli

import (
	"strings"
	"time"
)

// clock stores the current clock step and text.
type clock struct {
	timeFormat string
	clockRunes []string
}

// displayTime defines the visible values for a time.
type displayTime struct {
	time.Time
	digital string
	analog  string
}

// nolint
// unicode art clock spinners
const (
	asciiClock     = "/:﹣:\\:|"
	clockClock     = "🕛:🕐:🕑:🕒:🕓:🕔:🕕:🕖:🕗:🕘:🕙:🕚"
	brailleClock   = "⢎⡰:⢎⡡:⢎⡑:⢎⠱:⠎⡱:⢊⡱:⢌⡱:⢆⡱"
	brailleSpinner = " ⠁: ⠑: ⠰: ⡰:⢀⡠:⢄⡠:⢆⡀:⢎⡀:⢎ :⠎ :⠊ :⠈ "
)

func Clock(runes string) clock {
	return clock{
		timeFormat: TimeFormatHuman,
		clockRunes: strings.Split(runes, ":"),
	}
}

// Chars returns the clock chars.
func (c *clock) Chars(tickInterval time.Duration) string {
	step := (time.Now().UnixNano() / int64(tickInterval)) % int64(len(c.clockRunes))
	return c.clockRunes[step%int64(len(c.clockRunes))]
}

// DisplayTime returns the current clock text.
func (c *clock) DisplayTime(tickInterval time.Duration) *displayTime {
	t := time.Now()
	return &displayTime{t, t.Format(c.timeFormat), c.Chars(tickInterval)}
}
