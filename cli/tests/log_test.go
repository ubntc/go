package tests

import (
	"errors"
	"io"
	"log"
	"os"
	"testing"

	zlog "github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"github.com/ubntc/go/cli/cli"
	"github.com/ubntc/go/cli/loggers/stdlogger"
	"github.com/ubntc/go/cli/loggers/zerologger"
)

func TestZeroLog(t *testing.T) {
	cli.SetupLogging(zerologger.Setup)
	s := Capture(os.Stderr, func() {
		zlog.Info().Str("name", "abc").Msg("test")
	})
	assert.Contains(t, s, "test", "test message must be logged")
	assert.Contains(t, s, "name=", "test message must be logged")
	assert.Contains(t, s, "abc", "test message must be logged")
}

func TestStdLog(t *testing.T) {
	cli.SetupLogging(stdlogger.Setup)
	s := Capture(os.Stderr, func() {
		log.Printf("test name=%s\n", "abc")
	})
	assert.Contains(t, s, "test", "test message must be logged")
	assert.Contains(t, s, "name=abc", "test message must be logged")
}

func TestBadLog(t *testing.T) {
	bad := func(out io.Writer, timeFormat string) error {
		return errors.New("bad logger")
	}
	s := Capture(os.Stderr, func() {
		cli.SetupLogging(bad)
	})
	assert.Contains(t, s, "failed to setup logger")
	assert.Contains(t, s, "bad logger")
}
