package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"cloud.google.com/go/pubsub"

	clicks "github.com/ubntc/go/batching/batbq/_examples/ps2bq/clicks"
)

// clickPublisher is a simple dummy publisher for creating test clicks.
type clickPublisher struct {
	sync.Mutex
	total      int
	numWorkers int
	topic      *pubsub.Topic
}

// Next generates a new message.
func (p *clickPublisher) Next() ([]byte, error) {
	p.Lock()
	c := clicks.Click{
		ID:     fmt.Sprintf("click%d", p.total),
		Origin: "ps2bq-demo",
		Time:   time.Now().UTC(),
	}
	p.total++
	p.Unlock()
	return json.Marshal(c)
}

// Publish sends the given message to PubSub.
func (p *clickPublisher) Publish(ctx context.Context, data []byte) error {
	_, err := p.topic.Publish(ctx, &pubsub.Message{Data: data}).Get(ctx)
	return err
}

// Start starts the publisher.
func (p *clickPublisher) Start(ctx context.Context, sendInterval time.Duration, workerNum int) {
	p.Lock()
	p.numWorkers++
	p.Unlock()

	ticker := time.NewTicker(sendInterval)
	logtick := time.NewTicker(time.Second)
	start := time.Now()
	defer ticker.Stop()
	defer logtick.Stop()
	log.Printf("starting worker (sendInterval=%.2fs)", float64(sendInterval)/float64(time.Second))
	for {
		select {
		case <-ctx.Done():
			return
		case <-logtick.C:
			if workerNum != 1 {
				continue
			}
			p.Lock()
			mps := float64(p.total) / time.Now().Sub(start).Seconds()
			log.Printf("publisherStats(total=%d, mps=%.2f, workers=%d)", p.total, mps, p.numWorkers)
			p.Unlock()
		case <-ticker.C:
			data, err := p.Next()
			if err != nil {
				log.Print(err)
				return
			}

			err = p.Publish(ctx, data)
			if err != nil {
				log.Print(err)
				return
			}
		}
	}
}

func main() {
	var (
		project      = flag.String("p", os.Getenv("GOOGLE_CLOUD_PROJECT"), "Project ID")
		topic        = flag.String("t", "clicks", "Subscription Name")
		sendInterval = flag.Duration("d", time.Millisecond, "duration between test demo messages")
		concurrency  = flag.Int("c", 100, "concurrency level")
	)
	flag.Parse()

	ctx := context.Background()
	psClient, err := pubsub.NewClient(ctx, *project)
	if err != nil {
		log.Fatal(err)
	}
	defer psClient.Close()

	pub := &clickPublisher{topic: psClient.Topic(*topic)}
	dur := *sendInterval * time.Duration(*concurrency)

	for i := 1; i <= *concurrency; i++ {
		go pub.Start(ctx, dur, i)
	}
	<-ctx.Done()
}
