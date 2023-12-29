package pebble

import (
	"encoding/binary"

	"github.com/ubntc/go/kstore/kschema"
	"github.com/ubntc/go/kstore/provider/api"
)

type Message struct{ kschema.Message }

// Offset extracts and returns the offset from a given storage key.
func Offset(storageKey []byte) int64 {
	if len(storageKey) < 8 {
		panic("storageKey is too short")
	}

	return int64(binary.BigEndian.Uint64(storageKey[:8]))
}

// StorageKey returns a []byte key used for storing messages ordered by offset.
func StorageKey(msg api.Message) []byte {
	// Create a byte slice of size 8 (int64 size)
	offsetBytes := make([]byte, 8)
	// Encode the offset as big-endian into the byte slice
	binary.BigEndian.PutUint64(offsetBytes, uint64(msg.Offset()))

	// Concatenate offsetBytes with other identifying information
	return append(offsetBytes, msg.Key()...)
}

func (m *Message) StorageValue() []byte {
	return m.MustEncode()
}

// ensure we implement the full interface
func init() { _ = api.Message(&Message{}) }
