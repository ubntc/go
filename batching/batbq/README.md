[![Go Report Card](https://goreportcard.com/badge/github.com/ubntc/go/batcher/batbq)](https://goreportcard.com/report/github.com/ubntc/go/batcher/batbq)
[![cover-badge](https://img.shields.io/badge/coverage-91%25-brightgreen.svg?longCache=true&style=flat)](Makefile#10)

# Batched BigQuery Inserter

[![Go-batching Logo](resources/go-batching-logo.svg)](https://github.com/ubntc/go/blob/master/batching/batbq)

This package implements batching of messages for the `bigquery.Inserter`.

## Usage

```golang
package main

import (
	"context"
	"os"
	"time"

	"cloud.google.com/go/bigquery"
	"github.com/ubntc/go/batching/batbq"
	custom "github.com/ubntc/go/batching/batbq/_examples/simple/dummy"
)

var schema bigquery.Schema

func init() { schema, _ = bigquery.InferSchema(custom.Message{}) }

// Msg wraps the received data and implements batbq.Message.
type Msg struct {
	m *custom.Message // custom type providing data values and confirmation handlers
}

func (msg *Msg) Ack()           { msg.m.ConfirmMessage() }
func (msg *Msg) Nack(err error) {}
func (msg *Msg) Data() *bigquery.StructSaver {
	return &bigquery.StructSaver{InsertID: msg.m.ID, Struct: msg.m, Schema: schema}
}

func main() {
	source := custom.NewSource("src_name") // custom data source

	ctx := context.Background()
	client, _ := bigquery.NewClient(ctx, os.Getenv("GOOGLE_CLOUD_PROJECT"))
	output := client.Dataset("tmp").Table("batbq").Inserter()

	cfg := batbq.BatcherConfig{
		Capacity:      100,
		FlushInterval: time.Second,
	}

	input := make(chan batbq.Message, cfg.Capacity)
	batcher := batbq.NewInsertBatcher(cfg)

	go func() {
		source.Receive(ctx, func(m *custom.Message) { input <- &Msg{m} })
		close(input)
	}()
	batcher.Process(ctx, input, output)
}
```

Also see the [PubSub to BigQuery](_examples/pubsub-to-bq/main.go) example.


## Batcher Design

The package provides an `InsertBatcher` that requires an `input <-chan batbq.Message` channel to
collect individual messages from a streaming data source as shown in the [examples](./_examples).
The `InsertBatcher` also requires a `Putter` that implements `Put(context.Context, interface{})`
as provided by the regular `bigquery.Inserter`.

The `Put` method of a `bigquery.Inserter` will treat the given data as `bigquery.ValueSaver` or a
compatible type. Therefore batbq calls `batbq.Message.Data()` on each passed `batbq.Message`, which
must return a `*bigquery.StructSaver`.

Setting up a batch pipeline requires the following steps.

1. Create a wrapping type that implements `batbq.Message` providing `Ack()`, `Nack(error)`,
   and `Data() *bigquery.StructSaver`.
2. Create a `Putter` to receive the batches from the `InsertBatcher`.
3. Create a `chan batbq.Message` channel to pass data to the `InsertBatcher`.
4. In a goroutine, receive and wrap messages from a data source and send them to the channel.
5. Start the batcher with it's `Process(context.Context, <-chan batbq.Message, Putter)` method.

For instance, if your data source is PubSub, first register a message handler using
`subscription.Receive(ctx, handler)` in a goroutine, with the `handler` wrapping the
`pubsub.Message` in a `batbq.Message` and sending it to the input channel.
Then start the batcher to receive and batch these messages. The batcher will stop if the context
is canceled or the input channel is closed; there is no "stop" method.
See the full [PubSub to BigQuery](_examples/pubsub-to-bq/main.go) example for more details and
options.

## Worker Scaling

Internally batbq uses one or more worker goroutines to process data from the input channel.
If the `Putter` (usually a `bigquery.Inserter`) is stalled, the workers will block.
The worker will also block if the message confirmation is stalled by unanswered calls to `Ack()`
or `Nack(error)` for the currently processed batch.

If `BatcherConfig.AutoScale` is `true` a pipeline with slow senders or receivers is automatically
given more workers to increase the concurrency level. This results in more batches being collected
and sent concurrently via `output.Put(ctx, batch)`. However, all workers share the same
`input <-chan batbq.Message` and the same `output Putter`. Both, data source and output, must be
concurrency-safe by supporting concurrent calls of `Ack()`, `Nack(error)`, and `Put(ctx, batch)`.

## Multi Batching

The package also provides a `MultiBatcher` that can be set up to batch data from multiple inputs
and outputs in parallel. Please consult the corresponding [test case](multibatcher_test.go) on how
to set it up.
