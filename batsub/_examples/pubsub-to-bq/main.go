package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/pubsub"
	"github.com/ubntc/go/batsub"
)

// Click describes the context of a click on an Ad.
type Click struct {
	ID     string `json:"id"`
	Origin string `json:"origin"`
}

var clickSchema bigquery.Schema

func init() {
	var err error
	clickSchema, err = bigquery.InferSchema(Click{})
	if err != nil {
		panic(err)
	}
	dump, _ := json.Marshal(clickSchema)
	log.Printf("using click schema:\n%s", dump)
}

func pullMessages(ctx context.Context, rec batsub.Receiver, f batsub.BatchFunc) error {
	capacity := 100
	interval := time.Second
	sub := batsub.NewBatchedSubscription(rec, capacity, interval)

	if err := sub.ReceiveBatch(ctx, f); err != nil {
		return err
	}
	return nil
}

func insertMessages(ctx context.Context, ins *bigquery.Inserter, messages []*pubsub.Message) error {
	clicks := make([]*bigquery.StructSaver, 0, len(messages))
	errors := make(map[string]error)

	// parse messages
	for _, m := range messages {
		var click Click
		if err := json.Unmarshal(m.Data, &click); err != nil {
			errors[m.ID] = fmt.Errorf("failed to unmarshal click data: %v", m.Data)
			continue
		}
		clicks = append(clicks, &bigquery.StructSaver{
			Schema: clickSchema, InsertID: m.ID, Struct: &click,
		})
	}

	// insert messages and collect errors
	err := ins.Put(ctx, clicks)
	if mult, ok := err.(bigquery.PutMultiError); ok {
		for _, rowErr := range mult {
			errors[rowErr.InsertID] = &rowErr
		}
		err = nil
	} else if err != nil {
		return err
	}

	acked := 0
	// ack inserted messages and let messages of failed inserts expire
	for _, m := range messages {
		if err := errors[m.ID]; err != nil {
			log.Printf("insert failed, error: %s, ID: %s, Data: %v", err.Error(), m.ID, m.Data)
			// Do not `Nack` the message, instead let it expire.
			// This avoids the message being resent immediately.
			continue
		}
		acked++
		m.Ack()
	}

	log.Printf("acked and inserted messages %d/%d", acked, len(messages))
	return nil
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
	ins := bqClient.Dataset(*ds).Table(*table).Inserter()

	// start receiving and inserting messages
	err = pullMessages(ctx, subscription, func(ctx context.Context, messages []*pubsub.Message) {
		err := insertMessages(ctx, ins, messages)
		exitOnErr(err)
	})
	exitOnErr(err)
}
