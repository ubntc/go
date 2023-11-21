package kstore

import (
	"github.com/ubntc/go/kstore/kstore/config"
	"github.com/ubntc/go/kstore/kstore/status"
)

// tableSchema idenfies and fully defines a kstore table.
type tableSchema struct {
	Name   string      `json:"name,omitempty"`
	Topic  string      `json:"topic,omitempty"`
	Schema FieldSchema `json:"schema,omitempty"`

	state status.TableState
}

func NewTableSchema(name string, fields ...Field) (*tableSchema, error) {
	s := &tableSchema{
		Name:   name,
		Schema: fields,
		state:  status.TableState{},
	}
	if err := s.Validate(); err != nil {
		return nil, err
	}
	return s, nil
}

func (t *tableSchema) Validate() error {
	if t.Name == "" {
		return ErrorEmptyTableName
	}
	return nil
}

func (t *tableSchema) GetTopic() string {
	if t.Topic == "" {
		return config.DefaultTopicPrefix + t.Name
	}
	return t.Topic
}

func (t *tableSchema) GetTable() string {
	return t.Name
}

type TableTopic interface {
	Topic() string
	Table() string
}
