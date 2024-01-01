package kafkago

import (
	"fmt"

	"github.com/segmentio/kafka-go"
	"github.com/ubntc/go/kstore/provider/api"
)

type Message struct{ kafka.Message }

func (m *Message) Offset() uint64 {
	if m.Message.Offset < 0 {
		panic("negative offset")
	}
	return uint64(m.Message.Offset)
}

func (m *Message) Key() []byte   { return m.Message.Key }
func (m *Message) Value() []byte { return m.Message.Value }
func (m *Message) Topic() string { return m.Message.Topic }
func (m *Message) String() string {
	if m == nil {
		return "kafka.Message(nil)"
	}
	return fmt.Sprintf("kafka.Message(%s, %s, %s)", string(m.Key()), string(m.Value()), m.Topic())
}

// ensure we implement the full interface
func init() { _ = api.Message(&Message{}) }
