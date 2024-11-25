package buffertest

import (
	"context"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"
)

func Run(cfg Config) (int, error) {
	err := runPipeline(cfg)
	numHandled := mxNumHandled.Get()

	log.Info().Interface("results", Map{
		"topic":    cfg.Topic.Topic,
		"handled":  numHandled,
		"expected": cfg.NumEvents,
		"graceful": err == nil,
	}).Msg("pipeline finished")

	return numHandled, err
}

func runPipeline(cfg Config) error {
	ctx, cancel := context.WithTimeout(context.Background(), cfg.PipelineTimeout)
	defer cancel()

	client, err := newConn(ctx, cfg.Writer.Brokers[0], cfg.Topic.Topic)
	if err != nil {
		return errors.Wrap(err, "failed to create connection Kafka")
	}
	defer client.Close()

	log.Info().Msg("create writer")
	w := kafka.NewWriter(cfg.Writer)
	defer w.Close()

	log.Info().Msg("create topics")
	if err := client.CreateTopics(cfg.Topic); err != nil {
		return errors.Wrap(err, "failed to create topics")
	}

	log.Info().Msg("create reader")
	r := kafka.NewReader(cfg.Reader)
	defer r.Close()
	go consumeEvents(ctx, r)

	if err := produceEvents(ctx, w, cfg.NumEvents, cfg.WriterTick); err != nil {
		return errors.Wrap(err, "failed to produce events")
	}

	// wait for termination condition
	err = waitForEvents(ctx, cfg.NumEvents, cfg.WaiterTick)
	log.Info().Msg("closing reader, writer, and connection")
	return err
}
