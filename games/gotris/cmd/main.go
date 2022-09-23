package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ubntc/go/games/gotris/gotris"
)

func echo(values ...string) {
	fmt.Println(strings.Join(values, " "))
}

func main() {
	echo("starting Gotris")
	loop()
	echo("Gotris stopped")
}

// loop is the main Game loop, managing user input
// and state changes in the game in a step by step way.
func loop() {
	// TODO: allow speedup of ticker on higher levels
	game := gotris.NewGame(gotris.TestRules)

	reader := NewInputReader(os.Stdin)
	defer reader.Close()

	ticker := time.NewTicker(game.StepDuration)
	defer ticker.Stop()

	for {
		select {
		case key, more := <-reader.C:
			if !more {
				echo("stopping game loop", key)
				return
			}
			echo("input", key)
		case <-ticker.C:
			game.Advance()
			echo("advanced game", "step", strconv.Itoa(game.Steps))
			if game.Steps > 10 {
				reader.Close()
			}
		}
	}
}

type InputReader struct {
	source *os.File
	C      chan string
	once   sync.Once
}

func NewInputReader(source *os.File) *InputReader {
	r := &InputReader{source, make(chan string), sync.Once{}}
	// TODO: start reading single key presses from stdin
	return r
}

func (r *InputReader) Close() {
	r.once.Do(func() { close(r.C) })
}
