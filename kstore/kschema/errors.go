package kschema

import "errors"

var (
	ErrorTooManyValues        = errors.New("row has more values than schema fields")
	ErrorInvalidFieldType     = errors.New("row value has invalid field type")
	ErrorUnsupportedFieldType = errors.New("unsupported field type")
	ErrorEmptyTableName       = errors.New("Table.Name must not be empty")
	ErrorEmptyTopicName       = errors.New("Table.Topic must not be empty")
)
