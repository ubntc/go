package cli

import (
	"fmt"
)

// Prompt sets the global prompt message displayed in the interactive log line.
func Prompt(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	t := GetTerm()
	if t.IsDebug() {
		_, _ = t.Write([]byte("\r" + msg + "\n"))
	}
	t.SetMessage(msg)
}

// PromptVerbose sets the global prompt message if in debug mode.
func PromptVerbose(format string, v ...interface{}) {
	if GetTerm().IsVerbose() {
		Prompt(format, v...)
	}
}
