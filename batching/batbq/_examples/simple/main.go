// This file provides the content for the README.md it must

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
	batcher := batbq.NewInsertBatcher(capacity, interval, workers)

	go func() {
		source.Receive(ctx, func(m *custom.Message) { input <- &Msg{m} })
		close(input)
	}()
	batcher.Process(ctx, input, output)
}
