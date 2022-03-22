package buffertest

import (
	"context"
	"encoding/json"
	"io"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func newConn(ctx context.Context, address, topic string) (*kafka.Conn, error) {
	return kafka.DialLeader(ctx, "tcp", address, topic, 0)
}

func newEvent(topic string) (*kafka.Message, error) {
	id := uuid.NewString()
	data := Map{
		"timestamp": timestamppb.Now(),
		"id":        id,
		"topic":     topic,
	}
	value, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return &kafka.Message{
		Key:   []byte(id),
		Value: value,
	}, nil
}

func handleEvent(msg kafka.Message) error {
	n := mxNumHandled.Inc()
	if n%10 == 1 || true {
		log.Info().Interface("handleEvent", Map{
			"msg":         string(msg.Value),
			"topic":       msg.Topic,
			"num_handled": n,
		}).Msg("handled event")
	}
	return nil
}

func consumeEvents(ctx context.Context, r *kafka.Reader) {
	for {
		msg, err := r.ReadMessage(ctx)
		if err == io.EOF || err == context.Canceled {
			log.Info().Msg("reader drained")
			return
		}
		if err != nil {
			log.Error().Err(err).Msg("failed to read message")
			return
		}
		handleEvent(msg)
	}
}

func produceEvents(ctx context.Context, w *kafka.Writer, numEvents int, tick time.Duration) error {
	var kv *kafka.Message
	var err error
	ticker := time.NewTicker(tick)
	defer ticker.Stop()
	for i := 0; i < numEvents; i++ {
		kv, err = newEvent(w.Topic)
		if err != nil {
			return err
		}
		err = w.WriteMessages(ctx, *kv)
		if err != nil {
			return err
		}
		<-ticker.C
	}
	log.Info().Interface("produceEvents", Map{
		"num_events": numEvents,
		"topic":      w.Topic,
	}).Msg("produced all events")
	return nil
}

func waitForEvents(ctx context.Context, numEvents int, tick time.Duration) error {
	ticker := time.NewTicker(tick)
	defer ticker.Stop()
	for {
		n := mxNumHandled.Get()
		finished := n >= numEvents
		log.Info().Interface("waitForEvents", Map{
			"expected": numEvents,
			"handled":  n,
			"finished": finished,
		}).Msg("waiting for events")
		if finished {
			return nil
		}
		select {
		case <-ticker.C:
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
