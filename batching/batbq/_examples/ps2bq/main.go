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

	clicks "github.com/ubntc/go/batching/batbq/_examples/ps2bq/clicks"
	"github.com/ubntc/go/batching/batbq/patcher"
)

// ClickMessage wraps a Click and a pubsub.Message.
type ClickMessage struct {
	clicks.Click
	m *pubsub.Message
}

// NewClickMessage returns a new ClickMessage.
func NewClickMessage(m *pubsub.Message) (*ClickMessage, error) {
	msg := ClickMessage{m: m}
	if err := json.Unmarshal(m.Data, &msg.Click); err != nil {
		return nil, err
	}
	return &msg, nil
}

// Ack acks the underlying pubsub.Message.
func (c *ClickMessage) Ack() { c.m.Ack() }

// Nack prints the error.
func (c *ClickMessage) Nack(err error) {
	if err != nil {
		log.Print(err)
	}
}

// Data returns the BigQuery data using the pubsub message ID as bigquery InsertID.
func (c *ClickMessage) Data() bigquery.ValueSaver {
	return &bigquery.StructSaver{Schema: clickSchema, InsertID: c.m.ID, Struct: &c.Click}
}

var clickSchema bigquery.Schema

func init() {
	var err error
	clickSchema, err = bigquery.InferSchema(clicks.Click{})
	exitOnErr(err)
	dump, err := json.Marshal(clickSchema)
	exitOnErr(err)
	log.Printf("using click schema:\n%s", dump)
}

func exitOnErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func receive(ctx context.Context, sub *pubsub.Subscription, input chan<- batbq.Message) error {
	nMax := sub.ReceiveSettings.MaxOutstandingMessages
	sub.ReceiveSettings.MaxOutstandingBytes = nMax * 1000
	log.Printf("reading from subscription %s with MaxOutstandingMessages=%d", sub.ID(), nMax)
	return sub.Receive(ctx, func(ctx context.Context, m *pubsub.Message) {
		msg, err := NewClickMessage(m)
		if err != nil {
			log.Print(err)
			return
		}
		input <- msg
	})
}

func main() {
	var (
		project = flag.String("project", os.Getenv("GOOGLE_CLOUD_PROJECT"), "Project ID")
		topic   = flag.String("topic", "clicks", "Subscription Name")
		subfix  = flag.String("subfix", "", "subscription suffix")
		ds      = flag.String("ds", "tmp", "Dataset Name")
		table   = flag.String("table", "clicks", "Table Name")
		dry     = flag.Bool("dry", false, "setup pipeline but do not run the batcher")
		stats   = flag.Bool("stats", false, "print metrics on the console")
		cap     = flag.Int("cap", 500, "batch capacity")
	)
	flag.Parse()

	// setup PubSub client
	ctx := context.Background()
	psClient, err := pubsub.NewClient(ctx, *project)
	exitOnErr(err)
	defer psClient.Close()

	// setup BigQuery client
	bqClient, err := bigquery.NewClient(ctx, *project)
	exitOnErr(err)
	defer bqClient.Close()

	cfg := batbq.BatcherConfig{
		Capacity:      *cap,
		FlushInterval: time.Second,
		WorkerConfig: batbq.WorkerConfig{
			AutoScale: true,
		},
	}.WithDefaults()

	tab := bqClient.Dataset(*ds).Table(*table)
	input := make(chan batbq.Message, cfg.Capacity)
	output := tab.Inserter()

	sub := psClient.Subscription(*topic + *subfix)
	sub.ReceiveSettings.MaxOutstandingMessages = cfg.Capacity * cfg.MaxWorkers

	batcher := batbq.NewInsertBatcher("clicks", cfg)

	if *stats {
		batcher.Metrics().Watch(ctx)
	}

	if *dry {
		return
	}

	exitOnErr(patcher.Patch(ctx, tab, clickSchema))

	go func() {
		exitOnErr(receive(ctx, sub, input))
	}()

	exitOnErr(batcher.Process(ctx, input, output))
}
