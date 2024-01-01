package pebble_test

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"math/big"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ubntc/go/kstore/kschema"
	"github.com/ubntc/go/kstore/provider/pebble"
)

var (
	Msg            = kschema.RawMessage
	empty          = []byte{}
	k, v           = []byte("k"), []byte("v")
	zero, one, two = []byte("0"), []byte("1"), []byte("2")
	bignum         = []byte("9223372036854775807")

	maxOffsetBytes        = []byte{127, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}
	maxOffset      uint64 = 1<<63 - 1
)

func TestMessageOffsets(t *testing.T) {
	assert.Equal(t, maxOffset, pebble.Offset(maxOffsetBytes))

	type Test struct {
		name       string
		msg        kschema.Message
		wantBytes  []byte
		wantOffset uint64
	}

	tests := []Test{
		{"empty", Msg("t", 0, empty, empty), []byte{0, 0, 0, 0, 0, 0, 0, 0}, 0},
		{"zero", Msg("t", 0, zero, v), []byte{0, 0, 0, 0, 0, 0, 0, 0, zero[0]}, 0},
		{"one", Msg("t", 0, one, v), []byte{0, 0, 0, 0, 0, 0, 0, 0, one[0]}, 0},
		{"two", Msg("t", 0, two, v), []byte{0, 0, 0, 0, 0, 0, 0, 0, two[0]}, 0},
		{"big", Msg("t", 0, bignum, v), append(bytes.Repeat([]byte{0}, 8), bignum...), 0},
		{"max-empty", Msg("t", maxOffset, empty, v), maxOffsetBytes, maxOffset},
		{"max-big", Msg("t", maxOffset, bignum, v), append(maxOffsetBytes, bignum...), maxOffset},
	}

	numSamples := 1000
	// generate messages with various offsets
	for i := 0; i < numSamples; i++ {
		name := fmt.Sprintf("offset-%d", i)

		offset := uint64(i)
		offsetBytes := pebble.OffsetBytes(offset)
		msg := Msg("t", offset, k, v)

		tests = append(tests, Test{name, msg, append(offsetBytes, k...), offset})
	}

	// generate more random test cases with various offsets and bytes
	for i := 0; i < numSamples; i++ {
		name := fmt.Sprintf("random-%d", i)

		size, err := rand.Int(rand.Reader, big.NewInt(1000))
		assert.NoError(t, err)
		assert.LessOrEqual(t, size.Int64(), int64(1000))

		keySize := int(size.Int64())
		key := make([]byte, keySize)
		n, err := rand.Read(key)
		assert.Equal(t, keySize, n)
		assert.NoError(t, err)

		offset, err := rand.Int(rand.Reader, big.NewInt(int64(maxOffset)))
		assert.NoError(t, err)
		assert.LessOrEqual(t, len(offset.Bytes()), 8)

		offsetBytes := pebble.OffsetBytes(offset.Uint64())

		msg := Msg("t", offset.Uint64(), key, v)
		tests = append(tests, Test{name, msg, append(offsetBytes, key...), offset.Uint64()})
	}

	for _, tt := range tests {
		sk := pebble.StorageKey(&tt.msg)
		assert.Equal(t, tt.wantBytes, sk, "wrong key", tt.name)
		offset := pebble.Offset(sk)
		assert.Equal(t, tt.wantOffset, offset, "wrong offset", tt.name)
	}
}

// TestOffsetConversion shows that it is safe to convert between int64 and uint64,
// but only if the values are above 0 and below the max value for int64.
func TestOffsetConversion(t *testing.T) {
	var maxi int64 = 1<<63 - 1
	var maxu uint64 = 1<<63 - 1
	assert.Equal(t, int64(maxi), int64(maxu))
	assert.Equal(t, uint64(maxi), uint64(maxu))
	assert.Equal(t, maxu, maxOffset)

	// A negative int64 value converted to an uint64 will roll over beyond the maxOffset.
	var badValue int64 = -3
	assert.GreaterOrEqual(t, uint64(badValue), maxOffset)
}

func TestMessageOrder(t *testing.T) {
	// generate offsets with various offsets
	offsets := make([]uint64, 0)
	offsetBytes := make([][]byte, 0)
	var i uint64 = 0
	for i < maxOffset {
		offsets = append(offsets, i)
		offsetBytes = append(offsetBytes, pebble.OffsetBytes(i))
		switch {
		case i < 10:
			i++
		default:
			i *= 2
		}
	}

	// ensure that high offset messages have been added
	assert.Less(t, i, maxOffset*2)
	assert.GreaterOrEqual(t, i, maxOffset)

	// sort the bytes
	sort.Slice(offsetBytes, func(i, j int) bool {
		return bytes.Compare(offsetBytes[i], offsetBytes[j]) < 0
	})
	// ensure that the bytes have the same order as the offsets
	for i := 0; i < len(offsets); i++ {
		assert.Equal(t, offsets[i], pebble.Offset(offsetBytes[i]))
	}
}
