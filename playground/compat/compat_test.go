package compat

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ubntc/go/playground/compat/message"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

var marshaller = protojson.MarshalOptions{
	UseEnumNumbers: true,
}

var unmarshaller = protojson.UnmarshalOptions{
	DiscardUnknown: true,
}

func TestUseEnumNumbers(t *testing.T) {
	assert := assert.New(t)

	t.Log("compat test")

	msg := &message.Msg{
		Type: message.Type_TYPE_TWO,
	}

	b, err := marshaller.Marshal(msg)
	assert.NoError(err)
	assert.True(msg.Type == 2)
	v := string(b)
	assert.Contains(v, `:2`)
	t.Log("marshalled payload", v)

	t.Log("simulating receiving a new unknown enum value")
	v = strings.ReplaceAll(v, ":2", ":3")
	assert.Contains(v, `:3`)

	b = []byte(v)
	err = unmarshaller.Unmarshal(b, msg)
	assert.NoError(err)
	v = string(b)
	t.Log("modifed + unmarshalled payload:", v)

	switch msg.Type {
	case 0, 1, 2:
		t.Log("known payload, type:", message.Type_name[int32(msg.Type)])
	case 3:
		t.Log("unknown payload, type:", msg.Type)
	}

	assert.NotPanics(
		func() {
			typeVal := message.Type(10)
			t.Log("Go allows setting arbitrary enum values outside of the proto spec unknown_val:", typeVal)
		},
	)
}

func TestDefaultMarshaller(t *testing.T) {
	assert := assert.New(t)
	msg := &message.Msg{Type: 2}
	b, err := protojson.Marshal(msg)
	assert.NoError(err)
	v := string(b)
	t.Log("default marshalling with known enum", v)
	assert.Contains(v, `:"TYPE_TWO"`)

	msg = &message.Msg{Type: 3}
	b, err = protojson.Marshal(msg)
	assert.NoError(err)
	v = string(b)
	t.Log("default marshalling with unknown enum", v)
	assert.Contains(v, `:3`)
}

func TestProtoMarshaller(t *testing.T) {
	assert := assert.New(t)
	msg := &message.Msg{Type: 2}
	b, err := proto.Marshal(msg)
	assert.NoError(err)
	v := string(b)
	assert.Equal(v, "\b\x02")

	msg = &message.Msg{Type: 3}
	b, err = proto.Marshal(msg)
	assert.NoError(err)
	v = string(b)
	assert.Equal(v, "\b\x03")
}

func TestProtoUnmarshaller(t *testing.T) {
	assert := assert.New(t)
	msg := &message.Msg{}
	err := proto.Unmarshal([]byte("\b\x02"), msg)
	assert.NoError(err)
	assert.True(msg.Type == 2)
	assert.Equal(msg.Type, message.Type_TYPE_TWO)

	msg = &message.Msg{}
	err = proto.Unmarshal([]byte("\b\x03"), msg)
	assert.NoError(err)
	assert.True(msg.Type == 3)
	assert.NotEqual(msg.Type, 3, "types of ints should mismatch")
}
