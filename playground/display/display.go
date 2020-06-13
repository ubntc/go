package main

import (
	"context"
	"flag"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

var clearString = "\r" + strings.Repeat(" ", 100) + "\r"

// demo runtime
var runTime = flag.Duration("t", 3*time.Second, "runtime of the demo")

// demo function to produce messages
func sendMsg(ctx context.Context, ch chan<- string) {
	time.Sleep(time.Duration(rand.Uint32()) % *runTime)
	select {
	case ch <- fmt.Sprintf(time.Now().Format("2006-01-02 15:04:05.999")):
	case <-ctx.Done():
	}
}

// MessageDisplay displays messages.
type MessageDisplay struct {
	Messages chan string
	message  string
}

// Display displays incoming messages until the closing of the context.
func (d *MessageDisplay) Display(ctx context.Context) {
	ticker := time.Tick(100 * time.Millisecond)
	for {
		select {
		case <-ticker:
			fmt.Print(clearString + d.message + " ")
		case d.message = <-d.Messages:
		case <-ctx.Done():
			return
		}
	}
}

// Close closes the display.
func (d *MessageDisplay) Close() {
	fmt.Print(clearString)
}

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
