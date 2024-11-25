package main

import (
	"context"
	"fmt"
	"time"

	"github.com/ubntc/go/playground/pubsub/pkg"
)

// Goal: Design a basic Pub/Sub system!
// Actors:
// * MessageStore/Broker. Stores and delivers messages to Subscribers
// * Publisher/Producer. Sends messages to the Broker to be distributed to the Subscribers
// * Subscriber/Consumer. Receives messages from the Broker.

// Timeline:
// * 00:05 had a running main() and some types that did not do much
// * 00:15 added some methods and TODOs to implement basic Pub/Sub (not functional)
// * 00:30 basic most features drafted and partially implemented (saw where it needs to go)
// * 00:45 running demo with 10 messages sent and received
// * 01:00 added async subscribe, added backfill
// * 01:10 added retention handling
// * 01:20 adding more docs, moved to `pkg`
// * 01:30 fixed a DATA RACE caused by bad channel handling, added this timeline
// * 01:40 added readme
// * 01:45 fixed another DATA RACE (consumer list was not protected)

func main() {
	c := pkg.MyConsumer{Id: "c1"}
	b := pkg.NewBroker()
	p := pkg.Producer{
		Broker: b,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	fmt.Println("Starting Broker")
	go b.Run(ctx)

	fmt.Println("Starting Producer")
	go p.Run(ctx)

	go func() {
		<-time.After(time.Second * 6)
		fmt.Println("Adding Consumer")
		b.Subscribe(&c)
	}()

	<-ctx.Done()

	fmt.Println("Demo stopped")
}
