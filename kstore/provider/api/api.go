// This package defines interface types for implementing different persistence backends.

package api

import (
	"context"
)

type (
	TopicErrors map[string]error
	LoggerFunc  func(string, ...interface{})

	// Message defines the common interface for persistence messages send to and received from
	// the storage backend.
	//
	// This interface is complementary to the `kschema.Message` struct and should be used to
	// wrap existing types such as kafka-go's `kafka.Message`. Wrapping custom message types as
	// `api.Message`, the storage provider implementation can avoid passing around key and value
	// bytes, and message metadata from the original structs.
	Message interface {
		Key() []byte
		Value() []byte
		Offset() int64
		Topic() string
		String() string
	}

	Reader interface {
		Commit(ctx context.Context, msg Message) error
		Read(ctx context.Context) (Message, error)
		Close() error
	}

	Writer interface {
		Write(ctx context.Context, topic string, msg ...Message) error
		Close() error
	}

	Client interface {
		// direct client actions

		CreateTopics(ctx context.Context, topics ...string) (TopicErrors, error)
		DeleteTopics(ctx context.Context, topics ...string) (TopicErrors, error)
		Write(ctx context.Context, topic string, msg ...Message) error
		Close() error

		// derivded clients

		NewReader(topic string) Reader
		NewWriter() Writer

		// convenience funcs

		SetLogger(fn LoggerFunc)
		GetLogger() LoggerFunc

		// An "already exists" error is returned by many APIs to indicate that a table or topic
		// already exists. This is often considered a non-error when topics are created proactively.
		IsExistsError(err error) bool
	}
)
