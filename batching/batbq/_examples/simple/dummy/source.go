package dummy

import (
	"context"
	"fmt"
	"log"
	"time"

	"cloud.google.com/go/bigquery"
	"github.com/ubntc/go/batching/batbq"
)

// Source is a dummy source.
type Source struct {
	name      string
	Messages  []Message
	SendDelay time.Duration
}

// NewSource returns a dummy source.
func NewSource(name string) *Source {
	return &Source{name: name}
}

// Receive produces 200 raw dummy messages and sends them via the given message handler `f`.
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

// Chan returns a channel to receive the stored dummy messages as batbq.Message.
func (src *Source) Chan() <-chan batbq.Message {
	ch := make(chan batbq.Message, 100)
	schema, _ := bigquery.InferSchema(Message{})
	schema = schema.Relax()
	go func() {
		for _, m := range src.Messages {
			ch <- &batbq.LogMessage{
				bigquery.StructSaver{InsertID: m.ID, Struct: m, Schema: schema},
			}
			time.Sleep(src.SendDelay)
		}
		close(ch)
		// log.Print("send chan closed")
	}()
	return ch
}
