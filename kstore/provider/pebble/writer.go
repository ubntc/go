package pebble

import (
	"context"

	"github.com/ubntc/go/kstore/provider/api"
)

type Writer struct {
	client *Client
}

func NewWriter(c *Client) *Writer {
	return &Writer{client: c}
}

func (w *Writer) Write(ctx context.Context, topic string, messages ...api.Message) error {
	err := w.client.Write(ctx, topic, messages...)
	if err != nil {
		return err
	}
	Metrics.ObserveWrite(topic)
	return nil
}

func (w *Writer) Close() error {
	return nil
}

// ensure we implement the full interface
func init() { _ = api.Writer(&Writer{}) }
