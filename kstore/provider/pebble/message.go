package pebble

import (
	"math/big"

	"github.com/ubntc/go/kstore/kschema"
	"github.com/ubntc/go/kstore/provider/api"
)

type Message struct{ kschema.Message }

// Offset extracts and returns the offset from a given storage key.
func Offset(storageKey []byte) uint64 {
	if len(storageKey) < 8 {
		panic("storageKey is too short")
	}

	return big.NewInt(0).SetBytes(storageKey[:8]).Uint64()
}

// OffsetBytes converts an offset to a byte slice of 8 bytes.
func OffsetBytes(offset uint64) []byte {
	b := make([]byte, 8)
	offsetBytes := big.NewInt(0).SetUint64(offset).Bytes()
	copy(b[8-len(offsetBytes):], offsetBytes)
	return b
}

// StorageKey returns a []byte key used for storing messages ordered by offset.
func StorageKey(msg api.Message) []byte {
	// Concatenate offsetBytes with other identifying information
	return append(OffsetBytes(msg.Offset()), msg.Key()...)
}

func StorageValue(msg api.Message) []byte {
	m := Message{kschema.CopyMessage(msg)}
	return m.MustEncode()
}

// ensure we implement the full interface
func init() { _ = api.Message(&Message{}) }
