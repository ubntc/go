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
	"github.com/ubntc/go/metrics"

	clicks "github.com/ubntc/go/batching/batbq/_examples/ps2bq/clicks"
	"github.com/ubntc/go/batching/batbq/config"
	"github.com/ubntc/go/batching/batbq/multibatcher"
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
		project     = flag.String("project", os.Getenv("GOOGLE_CLOUD_PROJECT"), "Project ID")
		topic       = flag.String("topic", "clicks", "Subscription Name")
		subfix      = flag.String("subfix", "", "subscription suffix")
		ds          = flag.String("ds", "tmp", "Dataset Name")
		table       = flag.String("table", "clicks", "Table Name")
		dry         = flag.Bool("dry", false, "setup pipeline but do not run the batcher")
		stats       = flag.Bool("stats", false, "print metrics on the console")
		cap         = flag.Int("cap", 1000, "batch capacity")
		concurrency = flag.Int("c", 1, "number of independent batchers")
	)
	flag.Parse()

	ctx := context.Background()
	cfg := config.BatcherConfig{
		Capacity:      *cap,
		FlushInterval: time.Second,
		WorkerConfig: config.WorkerConfig{
			AutoScale: true,
		},
	}.WithDefaults()

	getInput := func(id string) <-chan batbq.Message {
		c, err := pubsub.NewClient(ctx, *project)
		exitOnErr(err)
		input := make(chan batbq.Message, cfg.Capacity)
		sub := c.Subscription(*topic + *subfix)
		sub.ReceiveSettings.MaxOutstandingMessages = 10000
		sub.ReceiveSettings.NumGoroutines = 10
		go func() {
			defer close(input)
			exitOnErr(receive(ctx, sub, input))
		}()
		return input
	}

	getOutput := func(id string) batbq.Putter {
		c, err := bigquery.NewClient(ctx, *project)
		exitOnErr(err)
		tab := c.Dataset(*ds).Table(*table)
		return tab.Inserter()
	}

	patch := func() error {
		c, err := bigquery.NewClient(ctx, *project)
		exitOnErr(err)
		t := c.Dataset(*ds).Table(*table)
		return patcher.PatchTable(ctx, t, clickSchema)
	}

	var batcherIDs []string
	for i := 1; i <= *concurrency; i++ {
		// all batchers use the same id to report to the same metrics
		batcherIDs = append(batcherIDs, "click")
	}

	mb := multibatcher.NewMultiBatcher(batcherIDs, batbq.Config(cfg))

	if *stats {
		metrics.Watch(ctx, mb.Metrics)
	}

	if *dry {
		return
	}

	exitOnErr(patch())
	exitOnErr(mb.MustProcess(ctx, getInput, getOutput))
}
