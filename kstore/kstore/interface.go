package kstore

import "context"

type (
	TopicErrors map[string]error
	LoggerFunc  func(string, ...interface{})

	Reader interface {
		Close() error
		Commit(ctx context.Context, msg Message) error
		Read(ctx context.Context) (Message, error)
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
		IsExistsError(err error) bool
	}
)
