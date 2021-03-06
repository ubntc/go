package zerologger

import (
	"io"
	"strings"

	stdlog "log"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type stdWrap struct {
	zerolog.Logger
}

func (w *stdWrap) Write(b []byte) (int, error) {
	w.Logger.Info().Msg(strings.TrimSpace(string(b)))
	return len(b), nil
}

// Setup sets up a zerolog.ConsoleWriter.
func Setup(out io.Writer, timeFormat string) error {
	logger := log.Output(zerolog.ConsoleWriter{
		Out:        out,
		TimeFormat: timeFormat,
	})
	log.Logger = logger

	slog := logger.With().Str("logger", "stdlog").Logger()
	stdlog.SetFlags(0)
	stdlog.SetOutput(&stdWrap{slog})

	log.Print("using zerolog.ConsoleWriter")
	return nil
}
