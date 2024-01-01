package pebble

import "errors"

var (
	ErrorTableNotInitalized = errors.New("pebble.DB not initalized")
	ErrorPipeClosed         = errors.New("pebble.Reader pipe closed")
	ErrorInvalidStartOffset = errors.New("pebble.Reader invalid start offset")
	ErrorOffsetNotFound     = errors.New("pebble.Reader could not find start offset")
	ErrorReicevedOldMessage = errors.New("pebble.Reader received old message from pebble.Client")
)
