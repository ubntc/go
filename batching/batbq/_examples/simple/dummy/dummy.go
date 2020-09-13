package dummy

import (
	"context"
	"fmt"
	"log"
	"time"
)

// Message is a dummy message.
type Message struct {
	ID  string
	Val int
}

// ConfirmMessage does nothing.
func (m *Message) ConfirmMessage() {}

// Source is a dummy source.
type Source struct {
	name string
}

// NewSource returns a dummy source.
func NewSource(name string) *Source {
	return &Source{name}
}

// Receive produces dummy messages.
func (src *Source) Receive(ctx context.Context, f func(m *Message)) {
	i := 0
	for {
		id := fmt.Sprint(time.Now().UnixNano())
		val := i
		select {
		case <-ctx.Done():
			return
		default:
			f(&Message{id, val})
			i++
		}
		if i >= 200 {
			log.Print("dummy source stopped")
			return
		}
		time.Sleep(time.Microsecond)
	}
}
