package zerologger

import (
	"io"
	"strings"

	stdlog "log"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const TimeformatHuman = "2006-01-02 15:04:05"

type stdWrap struct {
	zerolog.Logger
}

func (w *stdWrap) Write(b []byte) (int, error) {
	w.Logger.Info().Msg(strings.ReplaceAll(strings.ReplaceAll(strings.TrimSpace(string(b)), "\n", " "), "\r", ""))
	return len(b), nil
}

// Setup sets up a zerolog.ConsoleWriter.
func Setup(out io.Writer, timeFormat string) error {
	if timeFormat == "" {
		timeFormat = TimeformatHuman
	}
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
