package main

import (
	"flag"
	"fmt"
	"math/rand"
	"time"
)

var bogus = flag.Bool("bogus", false, "use bogus code")

func pause() {
	time.Sleep(time.Duration(rand.Uint32()%100) * time.Millisecond)
}

// This variant produces concurrent writes that can be detected via `go run --race`.
func bad() {
	stop := time.After(100 * time.Millisecond)
	var name string

	// start some producers doing concurrent writes (DANGER!)
	for i := 0; i < 10; i++ {
		go func(i int) {
			pause()
			name = fmt.Sprintf("name = %d", i)
		}(i)
	}

	// start consumer that shows the current value every 10ms
	go func() {
		tick := time.Tick(10 * time.Millisecond)
		for {
			select {
			case <-stop:
				return
			case <-tick:
				fmt.Println("read:", name)
			}
		}
	}()

	<-stop
}

// This variant is race-free.
func good() {
	stop := time.After(100 * time.Millisecond)
	names := make(chan string, 10)

	// start some producers concurrently writing to a channel (GOOD!)
	for i := 0; i < 10; i++ {
		go func(i int) {
			pause()
			select {
			case names <- fmt.Sprintf("name = %d", i):
			case <-stop:
			}
		}(i)
	}

	// start consumer that shows the current value every 10ms
	go func() {
		tick := time.Tick(10 * time.Millisecond)
		var name string
		for {
			select {
			case name = <-names:
			case <-stop:
				return
			case <-tick:
				fmt.Println("read:", name)
			}
		}
	}()

	<-stop
}

func main() {
	flag.Parse()
	if *bogus {
		bad()
	} else {
		good()
	}
}
