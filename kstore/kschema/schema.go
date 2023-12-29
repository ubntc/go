package kschema

import (
	"github.com/ubntc/go/kstore/kstore/config"
	"github.com/ubntc/go/kstore/kstore/status"
)

// Schema idenfies and fully defines a kstore table.
type Schema struct {
	Name   string      `json:"name,omitempty"`
	Topic  string      `json:"topic,omitempty"`
	Schema FieldSchema `json:"schema,omitempty"`

	state status.TableState
}

func NewTableSchema(name string, fields ...Field) (*Schema, error) {
	s := &Schema{
		Name:   name,
		Schema: fields,
		state:  status.TableState{},
	}
	if err := s.Validate(); err != nil {
		return nil, err
	}
	return s, nil
}

func (t *Schema) Validate() error {
	if t.Name == "" {
		return ErrorEmptyTableName
	}
	return nil
}

func (t *Schema) GetTopic() string {
	if t.Topic == "" {
		return config.DefaultTopicPrefix + t.Name
	}
	return t.Topic
}

func (t *Schema) GetTable() string {
	return t.Name
}

type TableTopic interface {
	Topic() string
	Table() string
}
