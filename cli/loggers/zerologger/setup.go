package zerologger

import (
	"io"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Setup sets up a zerolog.ConsoleWriter.
func Setup(out io.Writer, timeFormat string) error {
	logger := log.Output(zerolog.ConsoleWriter{
		Out:        out,
		TimeFormat: timeFormat,
	})
	log.Logger = logger

	log.Print("using zerolog.ConsoleWriter")
	return nil
}
