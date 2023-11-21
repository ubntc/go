package kafkago

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
	"github.com/ubntc/go/kstore/kstore"
)

type Writer struct {
	writer *kafka.Writer
}

func (w *Writer) Write(ctx context.Context, topic string, messages ...kstore.Message) error {
	km := make([]kafka.Message, 0)
	for _, m := range messages {
		// if i == 0 {
		// 	log.Printf("Writer.Write: %v to topic: %s (first of %d messages)", m.String(), topic, len(messages))
		// }
		km = append(km, kafka.Message{
			Topic: topic,
			Key:   m.Key(),
			Value: m.Value(),
		})
	}
	err := w.writer.WriteMessages(ctx, km...)
	if err != nil {
		log.Println(err)
	}
	return err
}

func (w *Writer) Close() error {
	return w.writer.Close()
}

// ensure we implement the full interface
func init() { _ = kstore.Writer(&Writer{}) }
