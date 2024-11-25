package pkg

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Broker is a message store and message distributer with a single topic.
// TODO: implement topics
type Broker struct {
	consumers []Consumer
	messages  []string
	// offsets  map[string]int64 // TODO: implement offset tracking per consumer

	mu          sync.Mutex
	newMessage  chan string
	newConsumer chan Consumer
}

func NewBroker() *Broker {
	return &Broker{
		newMessage:  make(chan string, 100),
		newConsumer: make(chan Consumer),
	}
}

// Send is used by message publishers to add new messages on the broker.
func (b *Broker) Send(msg string) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.messages = append(b.messages, msg)
	fmt.Printf("added message %s, len = %d\n", msg, len(b.messages))
	b.newMessage <- msg
	return nil
}

// Subscribe allows a Consumer to subscribe to the broker messages.
// The broker will be also issue a backfill of the consumer, sending it historical messages.
func (b *Broker) Subscribe(consumer Consumer) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.consumers = append(b.consumers, consumer)
	b.newConsumer <- consumer
}

// Run starts the broker.
// Run is a blocking call and should be started as goroutine.
func (b *Broker) Run(ctx context.Context) {
	cleaner := time.NewTicker(time.Second * 3)
	for {
		select {
		case <-ctx.Done():
			return

		case msg := <-b.newMessage:
			b.mu.Lock()
			for _, cons := range b.consumers {
				// TODO: keep track of offsets per consumer
				cons.Receive(msg)
			}
			b.mu.Unlock()

		case c := <-b.newConsumer:
			b.mu.Lock()
			fmt.Println("backfilling new consumer")
			for _, msg := range b.messages {
				c.Receive(msg)
			}
			fmt.Println("backfill finished")
			b.mu.Unlock()

		case <-cleaner.C:
			b.mu.Lock()
			if len(b.messages) > 0 {
				fmt.Println("cleaning up messages (strategy drop 50%)")
				b.messages = b.messages[len(b.messages)/2:]
			}
			b.mu.Unlock()
		}
	}
}
