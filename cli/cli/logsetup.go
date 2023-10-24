package cli

import (
	"fmt"
	"io"
	"os"
	"time"
)

// Supported logger types
const (
	// TimeFormatHuman is the only reasonable time format for human readability.
	TimeFormatHuman = time.DateTime
	// ClearAll is a fallback clear string.
	ClearAll = "\r                                                                                          \r"
)

// LogSetupFunc configures a specific Logger.
type LogSetupFunc func(out io.Writer, timeFormat string) error

// stdErrWrapper wraps stderr to allow replacing stdErr by other libraries and in testing.
type stdErrWrapper struct{}

func (w *stdErrWrapper) Write(b []byte) (int, error) {
	return os.Stderr.Write(b)
}

// SetupLogging configures interactive logging for commandline applications.
func SetupLogging(setup LogSetupFunc) {
	// Wrap stderr to avoid clock chars to mix with actual log lines.
	// The term will then become the new io.Writer for the logger.
	// Writing log lines to the term first clears the interactive line
	// before the actual log line is written.
	term := GetTerm()
	term.SetOutput(&stdErrWrapper{})

	if err := setup(term, TimeFormatHuman); err != nil {
		os.Stderr.WriteString(fmt.Sprintf("failed to setup logger, error: %v\n", err))
		os.Stderr.Sync()
	}
}
