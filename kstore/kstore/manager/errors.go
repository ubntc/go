package manager

import "errors"

var (
	ErrorWriterNotDefined = errors.New("Writer not defined")
	ErrorEmptyTopic       = errors.New("SchemaManager.Topic not set")
)
