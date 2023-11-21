package kstore

import "fmt"

// Message defines the accepted interface for persistence messages
// send to and received from the cloud storage backend.
type Message interface {
	Key() []byte
	Value() []byte
	String() string
}

// message is the internal message type used to generate concrete messages.
type message struct {
	key   []byte
	value []byte
}

func (m *message) Key() []byte   { return m.key }
func (m *message) Value() []byte { return m.value }
func (m *message) String() string {
	if m == nil {
		return "kstore.Message(nil)"
	}
	return fmt.Sprintf("kstore.Message(%s, %s)", string(m.key), string(m.value))
}
