// This file provides the content for the README.md it must

package main

import (
	"context"
	"flag"
	"log"
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

// Ack acknowledges the original message.
func (msg *Msg) Ack() { msg.m.ConfirmMessage() }

// Nack handles insert errors.
func (msg *Msg) Nack(err error) {
	if err != nil {
		log.Print(err)
	}
}

// Data returns the message as bigquery.StructSaver.
func (msg *Msg) Data() bigquery.ValueSaver {
	return &bigquery.StructSaver{InsertID: msg.m.ID, Struct: msg.m, Schema: schema}
}

var dry = flag.Bool("dry", false, "setup pipeline but do not run the batcher")

func main() {
	flag.Parse()
	source := custom.NewSource("src_name") // custom data source

	ctx := context.Background()
	client, _ := bigquery.NewClient(ctx, os.Getenv("GOOGLE_CLOUD_PROJECT"))
	output := client.Dataset("tmp").Table("batbq").Inserter()

	cfg := batbq.BatcherConfig{
		Capacity:      100,
		FlushInterval: time.Second,
	}

	input := make(chan batbq.Message, cfg.Capacity)
	batcher := batbq.NewInsertBatcher("custom_message", cfg)

	go func() {
		source.Receive(ctx, func(m *custom.Message) { input <- &Msg{m} })
		close(input)
	}()

	if *dry {
		return
	}
	batcher.Process(ctx, input, output)
}
