package main

import (
	"context"
	"flag"
	"fmt"
	"math/rand/v2"
	"strings"
	"time"
)

// MessageDisplay displays messages.
type MessageDisplay struct {
	Messages chan string

	nextMsg    string
	currentMsg string
}

// Display displays incoming messages until the closing of the context.
func (d *MessageDisplay) Display(ctx context.Context) {
	ticker := time.Tick(100 * time.Millisecond)
	for {
		select {
		case <-ticker:
			fmt.Print(d.ClearString() + d.nextMsg)
		case d.nextMsg = <-d.Messages:
		case <-ctx.Done():
			return
		}
	}
}

// Close clears and closes the display.
func (d *MessageDisplay) Close() {
	fmt.Print(d.ClearString())
}

// ClearString returns a string that clears the current message without creating a new line.
func (d *MessageDisplay) ClearString() string {
	return "\r" + strings.Repeat(" ", len(d.currentMsg)) + "\r"
}

// demo function to produce messages
func sendMsg(ctx context.Context, ch chan<- string) {
	time.Sleep(time.Duration(rand.Uint32()) % *runTime)
	select {
	case ch <- time.Now().Format(time.DateTime + ".999"):
	case <-ctx.Done():
	}
}

var runTime = flag.Duration("t", 3*time.Second, "runtime of the demo")

func main() {
	flag.Parse()
	d := &MessageDisplay{Messages: make(chan string)}
	ctx, cancel := context.WithCancel(context.Background())

	go d.Display(ctx)
	defer d.Close()

	for i := 0; i < 1000; i++ {
		go sendMsg(ctx, d.Messages)
	}

	time.Sleep(*runTime)
	cancel()
	<-ctx.Done()
}
