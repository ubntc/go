// This package defines a common message types use by kstore.

package kschema

import (
	"encoding/json"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/ubntc/go/kstore/provider/api"
)

var globalOffset atomic.Uint64

func init() { globalOffset.Add(uint64(time.Now().UnixNano())) }

// fields defines a marshallable message
type fields struct {
	Topic  string `json:"topic,omitempty"`
	Key    []byte `json:"key,omitempty"`
	Value  []byte `json:"value,omitempty"`
	Offset uint64 `json:"offset,omitempty"`
}

// Message is the defaulr message type used to create new concrete payloads
// and implements the api.Message interface.
type Message struct{ fields fields }

func (m *Message) Encode() ([]byte, error) {
	return json.Marshal(m.fields)
}

func (m *Message) Decode(data []byte) error {
	if m == nil {
		*m = Message{}
	}
	return json.Unmarshal(data, &m.fields)
}

func (m *Message) MustEncode() []byte {
	data, err := m.Encode()
	if err != nil {
		panic(err)
	}
	return data
}

func NewMessage(topic string, key, value []byte) *Message {
	offset := globalOffset.Add(1)
	msg := RawMessage(topic, offset, key, value)
	return &msg
}

func CopyMessage(msg api.Message) Message {
	return Message{fields{Topic: msg.Topic(), Offset: msg.Offset(), Key: msg.Key(), Value: msg.Value()}}
}

func RawMessage(topic string, offset uint64, key, value []byte) Message {
	return Message{fields{Topic: topic, Offset: offset, Key: key, Value: value}}
}

func (m *Message) Key() []byte    { return m.fields.Key }
func (m *Message) Value() []byte  { return m.fields.Value }
func (m *Message) Offset() uint64 { return m.fields.Offset }
func (m *Message) Topic() string  { return m.fields.Topic }
func (m *Message) String() string {
	if m == nil {
		return "api.Message(nil)"
	}
	return fmt.Sprintf("api.Message(%s, %s, %s, %d)", m.Topic(), string(m.Key()), string(m.Value()), m.Offset())
}

var _ = api.Message(&Message{})
