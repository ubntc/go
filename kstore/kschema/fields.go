package kschema

import (
	"errors"

	"github.com/ubntc/go/kstore/provider/api"
)

type FieldType string

const (
	FieldTypeString  = "string"
	FieldTypeInt64   = "int64"
	FieldTypeFloat64 = "float64"
	FieldTypeBool    = "bool"
	FieldTypeRecord  = "record"
)

type Field struct {
	Name string    `json:"name,omitempty"`
	Type FieldType `json:"type,omitempty"`
	// Repeated   bool   `json:"repeated,omitempty"`
	// RecordType Schema `json:"record_type,omitempty"`
}

type FieldSchema []Field

func (s FieldSchema) Validate(row Row) error {
	return s.ValidateRows(row)
}

func (s FieldSchema) ValidateRows(rows ...Row) error {
	var err error
	if s == nil {
		// allow schemaless KV tables
		return nil
	}
	if len(rows) == 0 {
		return nil
	}
	for _, row := range rows {
		for i, v := range row.GetValues() {
			if i >= len(s) {
				err = errors.Join(err, ErrorTooManyValues)
				break
			}
			field := s[i]
			switch field.Type {
			case FieldTypeString:
				if _, ok := v.(string); !ok {
					err = errors.Join(err, ErrorInvalidFieldType)
				}
			default:
				err = errors.Join(err, ErrorUnsupportedFieldType)
			}
		}
		// return all errors from the first erroneous row
		if err != nil {
			return err
		}
	}
	return nil
}

func (s FieldSchema) ValidateMessage(msg api.Message) error {
	if msg == nil {
		return nil
	}
	row := Row{}
	if err := row.Decode(msg.Value()); err != nil {
		return err
	}
	return s.ValidateRows(row)
}
