package pebble

import (
	"context"

	"github.com/cockroachdb/pebble"
	"github.com/ubntc/go/kstore/kschema"
	"github.com/ubntc/go/kstore/provider/api"
)

type Writer struct {
	client *Client
}

func NewWriter(c *Client) *Writer {
	return &Writer{client: c}
}

func (w *Writer) Write(ctx context.Context, topic string, messages ...api.Message) error {
	db, err := w.client.GetDB(topic)
	if err != nil {
		return err
	}

	for _, m := range messages {
		msg := Message{kschema.CopyMessage(m)}
		if err := db.Set(StorageKey(&msg), msg.StorageValue(), pebble.Sync); err != nil {
			return err
		}
	}
	return nil
}

func (w *Writer) Close() error {
	return nil
}

// ensure we implement the full interface
func init() { _ = api.Writer(&Writer{}) }
