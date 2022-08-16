package pkg

import (
	"context"
	"fmt"
	"time"
)

type Producer struct {
	// Broker is a reference to the broker where to publish new messages
	Broker *Broker
}

// Run is the main event loop of a Producer.
// In this loop it watches for system changes and sends the corresponding events to the broker.
// Run is a blocking call and should be started as goroutine.
func (p *Producer) Run(ctx context.Context) {
	i := 0
	ticker := time.NewTicker(time.Second)
	for {
		i += 1
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			p.Broker.Send(fmt.Sprintf("message %d", i))
		}
	}
}
