package kafkago

import (
	"context"

	"github.com/segmentio/kafka-go"
	"github.com/ubntc/go/kstore/provider/api"
)

type Reader struct {
	topic  string
	reader *kafka.Reader
}

func (r *Reader) Close() error {
	return r.reader.Close()
}

func (r *Reader) Commit(ctx context.Context, msg api.Message) error {
	return r.reader.CommitMessages(ctx, kafka.Message{
		Topic: r.topic,
		Key:   msg.Key(),
		Value: msg.Value(),
	})
}

func (r *Reader) Read(ctx context.Context) (api.Message, error) {
	m, err := r.reader.FetchMessage(ctx)
	if err != nil {
		return nil, err
	}
	return &Message{m}, nil
}

// ensure we implement the full interface
func init() { _ = api.Reader(&Reader{}) }
