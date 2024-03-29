package kstore

import (
	"context"
	"errors"
	"io"
	"log"

	"github.com/ubntc/go/kstore/provider/api"
)

var (

	// Manager Validation

	ErrorNotInitialized      = errors.New("SchemaManager not initialized")
	ErrorWriterTopicNotEmpty = errors.New("Writer.Topic must not be set")

	// Data Validation

	ErrorNoTableSchema = errors.New("No schema found for table")

	ErrorTopicMismatch = errors.New("Message topic and Table topic do not match")
	ErrorNilRow        = errors.New("Row was nil while setting key or value")
	ErrorNilMessage    = errors.New("Message was nil while setting key or value")

	// Store Validation

	ErrorReadStoreNotInitalized = errors.New("Store not initialized before reading")
	ErrorStoreNotInitalized     = errors.New("Store not initialized")
)

// ChanGo runs a function as goroutine and returns the returned error (or nil) on a non-blokcing error channel.
func ChanGo(fn func() error) <-chan error {
	errch := make(chan error, 1)
	go func() {
		errch <- fn()
	}()
	return errch
}

func FilterGraceful(err error) error {
	switch err {
	case context.Canceled, io.EOF:
		return nil
	case io.ErrClosedPipe, io.ErrUnexpectedEOF:
		// TODO: Find out if these are graceful or not, i.e.,
		//       What happens during a K8s deploy?
		//       What happens on SIGINT SIGTERM?
		//       How can a Kafka "close a pipe"?
		//       What happens during a network outage?
		//       Find more cases, and test them!
		return err
	default:
		return err
	}
}

func NewLogger(name string) api.LoggerFunc {
	return func(format string, args ...any) {
		log.Printf(name+": "+format, args...)
		log.Println()
	}
}

func NilLogger() api.LoggerFunc {
	return func(format string, args ...any) {}
}
