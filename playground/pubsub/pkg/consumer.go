package pkg

import "fmt"

// Consumer is the interface for registering consumers at the broker.
// The broker will call the consumer's Receive method to forward messages from the broker.
type Consumer interface {
	Receive(string)
}

// MyConsumer is an example implementation of a Consumer.
type MyConsumer struct {
	Id string
}

// Receive implements the Consumer interface
func (c *MyConsumer) Receive(message string) {
	fmt.Println("received message", message)
}
