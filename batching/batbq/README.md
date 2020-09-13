[![Go Report Card](https://goreportcard.com/badge/github.com/ubntc/go/batcher/batbq)](https://goreportcard.com/report/github.com/ubntc/go/batcher/batbq)
[![cover-badge](https://img.shields.io/badge/coverage-93%25-brightgreen.svg?longCache=true&style=flat)](Makefile#10)

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
	capacity, interval, workers := 100, time.Second, 1

	source := custom.NewSource("src_name") // custom data source

	ctx := context.Background()
	client, _ := bigquery.NewClient(ctx, os.Getenv("GOOGLE_CLOUD_PROJECT"))
	output := client.Dataset("tmp").Table("batbq").Inserter()

	input := make(chan batbq.Message, capacity)
	batcher := batbq.NewInsertBatcher(batbq.BatcherConfig{capacity, interval, workers, 0})

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
as provided the regular `bigquery.Inserter`.

The `Put` method of a `bigquery.Inserter` will treat the given data as `bigquery.ValueSaver` or a
compatible type. Therefore batbq calls `batbq.Message.Data()` on each passed `batbq.Message`, which
must return a `*bigquery.StructSaver`.

Setting up a batch pipeline requires the following steps.

1. Create a wrapping type that implements `batbq.Message` providing `Ack()`, `Nack(error)`,
   and `Data() *bigquery.StructSaver`.
2. Create a `chan batbq.Message` channel to pass data to the `InsertBatcher`.
3. Consume messages from the  original data source and fill them into the channel.

For instance, for PubSub you need to register a handler using `subscription.Receive(ctx, handler)`
and in the `handler` you need to convert the `pubsub.Message` to a `batbq.Message` as shown in the
[PubSub to BigQuery](_examples/pubsub-to-bq/main.go) example.

## Worker Scaling

Internally batbq uses one or more worker goroutines to process data from the input channel.
If the `Putter` (e.g., a `bigquery.Inserter`) is stalled, the workers will block.
The worker will also block if an inserted batch of messages is not yet confirmed on the sender side,
i.e., if one or more `Ack()` or `Nack(error)` calls are blocking.

If `BatcherConfig.AutoScale` is `true` a pipeline with slow senders or receivers is automatically
given more workers to increase the concurrency level. This results in more batches being collected
and sent concurrently via `output.Put(ctx, batch)`. However, all workers share the same
`input <-chan batbq.Message` and the same `output Putter`. Both the data source and the output must
be concurrency-safe and allow for concurrent calls of `Ack()`, `Nack(error)`, and `Put(ctx, batch)`.

## Multi Batching

The package also provides a `MultiBatcher` that can be set up to batch data from multiple inputs
and outputs in parallel. Please consult the corresponding [test case](multibatcher_test.go) on how
to set it up.
