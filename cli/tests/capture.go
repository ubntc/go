package tests

import (
	"bytes"
	"io"
	"os"
)

// Capture captures stderr or stdout.
func Capture(file *os.File, fn func()) string {
	var buf bytes.Buffer
	r, w, err := os.Pipe()
	if err != nil {
		panic(err)
	}

	stdout := os.Stdout
	stderr := os.Stderr

	switch file {
	case os.Stdout:
		defer func() { os.Stdout = stdout }()
		os.Stdout = w
	case os.Stderr:
		defer func() { os.Stderr = stderr }()
		os.Stderr = w
	default:
		panic("unsupported file descripitor" + file.Name())
	}

	fn()
	w.Close()
	io.Copy(&buf, r)
	return buf.String()
}
