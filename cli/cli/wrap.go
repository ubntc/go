package cli

import (
	"io"
	"os"
)

// CrWriter is an io.Writer that adds a CR character to every NL character written to it
type CrWriter struct {
	w    io.Writer
	last byte
}

// writeByte writes CR before it writes NL bytes
func (c *CrWriter) writeByte(b byte) error {
	var bytes []byte
	if b == '\n' && c.last != '\r' {
		bytes = []byte{'\r', b}
	} else {
		bytes = []byte{b}
	}
	if _, err := c.w.Write(bytes); err != nil {
		return err
	}
	c.last = b
	return nil
}

// Write implements io.Writer. NL bytes written get prepened CR bytes.
func (c *CrWriter) Write(p []byte) (n int, err error) {
	for _, b := range p {
		err = c.writeByte(b)
		if err != nil {
			return
		}
		n++ // extra bytes are ignored to not confuse the writers
	}
	return
}

// CrPipe returns a regular os.File for writing bytes. The pipe prepends CR bytes
// to non-prepended NL bytes and writes the result to the destination file.
func CrPipe(dst *os.File) (*os.File, error) {
	// 1. We create a CrWriter with the given os.File as final destination.
	cr := &CrWriter{dst, 0}

	// 2. To get a new writer that is a regular os.File, we need to use an os.Pipe.
	//    In an os.Pipe, the reads from the reader return the bytes written to the writer.
	r, exposedWriter, err := os.Pipe()
	if err != nil {
		return nil, err
	}

	// 3. To connect the returned writer to the CrWriter, we need to read from the os.Pipe.
	//    This can be done easily by reading from pipe's reader and writing to the CrWriter.
	go func() {
		defer r.Close()
		_, _ = io.Copy(cr, r)
	}()

	// 4. Return the new and fully connected writer as regular os.File.
	return exposedWriter, nil
}
