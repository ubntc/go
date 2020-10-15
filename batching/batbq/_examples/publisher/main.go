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
	total        int
	numWorkers   int
	project      string
	topic        string
	sendInterval time.Duration
}

// Next generates a new message.
func (p *clickPublisher) Next(prefix string) ([]byte, error) {
	p.Lock()
	num := p.total
	p.total++
	p.Unlock()

	c := clicks.Click{
		ID:     fmt.Sprintf("click_%s_%d", prefix, num),
		Origin: "ps2bq-demo",
		Time:   time.Now().UTC(),
	}
	return json.Marshal(c)
}

// Start starts the publisher.
func (p *clickPublisher) Worker(ctx context.Context) {
	p.Lock()
	p.numWorkers++
	workerNum := p.numWorkers
	p.Unlock()

	psClient, err := pubsub.NewClient(ctx, p.project)
	if err != nil {
		log.Fatal(err)
	}
	defer psClient.Close()

	prefix := fmt.Sprintf("pub_%d", workerNum)
	psTopic := psClient.Topic(p.topic)
	msgRate := float64(p.sendInterval) / float64(time.Second)

	ticker := time.NewTicker(p.sendInterval)
	logtick := time.NewTicker(time.Second)
	defer ticker.Stop()
	defer logtick.Stop()

	results := make(chan *pubsub.PublishResult, 100)
	defer close(results)

	go func() {
		for res := range results {
			if _, err := res.Get(ctx); err != nil {
				log.Print(err)
			}
		}
	}()

	log.Printf("starting worker with ticker duration %.2fs and expected message rate %.2f)",
		p.sendInterval.Seconds(), msgRate,
	)

	start := time.Now()
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
			data, err := p.Next(prefix)
			if err != nil {
				log.Print(err)
				return
			}
			results <- psTopic.Publish(ctx, &pubsub.Message{Data: data})
		}
	}
}

func main() {
	var (
		project      = flag.String("p", os.Getenv("GOOGLE_CLOUD_PROJECT"), "Project ID")
		topic        = flag.String("t", "clicks", "Subscription Name")
		sendInterval = flag.Duration("d", time.Millisecond, "duration between clicks")
		concurrency  = flag.Int("c", 1, "concurrency level")
	)
	flag.Parse()

	ctx := context.Background()
	dur := *sendInterval
	pub := &clickPublisher{
		topic:        *topic,
		project:      *project,
		sendInterval: dur,
	}

	for i := 1; i <= *concurrency; i++ {
		go pub.Worker(ctx)
	}
	<-ctx.Done()
}
