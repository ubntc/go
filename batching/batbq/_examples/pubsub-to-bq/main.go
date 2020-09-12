package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"os"
	"time"

	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/pubsub"
	"github.com/ubntc/go/batching/batbq"
)

// Click describes the context of a click on an Ad.
type Click struct {
	ID     string `json:"id"`
	Origin string `json:"origin"`
}

// ClickMessage wraps a Click and a pubsub.Message.
type ClickMessage struct {
	Click
	m *pubsub.Message
}

// NewClickMessage returns a new ClickMessage.
func NewClickMessage(m *pubsub.Message) (*ClickMessage, error) {
	var msg ClickMessage
	if err := json.Unmarshal(m.Data, &msg.Click); err != nil {
		return nil, err
	}
	return &msg, nil
}

// Ack acks the underlying pubsub.Message.
func (c *ClickMessage) Ack() { c.m.Ack() }

// Nack prints the error.
func (c *ClickMessage) Nack(err error) {
	log.Print(err)
}

// Data returns the BigQuery data using the pubsub message ID as bigquery InsertID.
func (c *ClickMessage) Data() *bigquery.StructSaver {
	return &bigquery.StructSaver{Schema: clickSchema, InsertID: c.m.ID, Struct: &c.Click}
}

var clickSchema bigquery.Schema

func init() {
	var err error
	clickSchema, err = bigquery.InferSchema(Click{})
	exitOnErr(err)
	dump, _ := json.Marshal(clickSchema)
	log.Printf("using click schema:\n%s", dump)
}

func exitOnErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	var (
		project = flag.String("project", os.Getenv("GOOGLE_CLOUD_PROJECT"), "Project ID")
		sub     = flag.String("sub", "clicks", "Subscription Name")
		ds      = flag.String("ds", "tmp", "Dataset Name")
		table   = flag.String("table", "clicks", "Table Name")
	)
	flag.Parse()

	// setup PubSub source
	ctx := context.Background()
	psClient, err := pubsub.NewClient(ctx, *project)
	exitOnErr(err)
	defer psClient.Close()
	subscription := psClient.Subscription(*sub)

	// setup BQ sink
	bqClient, err := bigquery.NewClient(ctx, *project)
	exitOnErr(err)
	defer bqClient.Close()
	output := bqClient.Dataset(*ds).Table(*table).Inserter()

	capacity := 10
	interval := time.Second
	workers := 1

	input := make(chan batbq.Message, capacity)
	batcher := batbq.NewInsertBatcher(batbq.BatcherConfig{capacity, interval, workers, 0})

	go subscription.Receive(ctx, func(ctx context.Context, m *pubsub.Message) {
		msg, err := NewClickMessage(m)
		if err != nil {
			log.Print(err)
		}
		input <- msg
	})
	batcher.Process(ctx, input, output)
}
