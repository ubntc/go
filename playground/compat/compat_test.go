package compat

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ubntc/go/playground/compat/message"
	"google.golang.org/protobuf/encoding/protojson"
)

var marshaller = protojson.MarshalOptions{
	UseEnumNumbers: true,
}

var unmarshaller = protojson.UnmarshalOptions{
	DiscardUnknown: true,
}

func TestUseEnumNumbers(t *testing.T) {
	assert := assert.New(t)

	fmt.Println("compat test")

	msg := &message.Msg{
		Type: message.Type_TYPE_TWO,
	}

	b, err := marshaller.Marshal(msg)
	assert.NoError(err)
	assert.True(msg.Type == 2)
	v := string(b)
	assert.Contains(v, `:2`)
	fmt.Println("marshalled payload", v)

	fmt.Println("simulating receiving a new unknown enum value")
	v = strings.ReplaceAll(v, ":2", ":3")
	assert.Contains(v, `:3`)

	b = []byte(v)
	err = unmarshaller.Unmarshal(b, msg)
	assert.NoError(err)
	v = string(b)
	fmt.Println("modifed + unmarshalled  payload:", v)

	switch msg.Type {
	case 0, 1, 2:
		fmt.Println("known payload, type:", message.Type_name[int32(msg.Type)])
	case 3:
		fmt.Println("unknown payload, type:", msg.Type)
	}

	assert.NotPanics(
		func() {
			typeVal := message.Type(10)
			fmt.Println("Go allows setting arbitrary enum values outside of the proto spec unknown_val:", typeVal)
		},
	)
}

func TestDefaultMarshaller(t *testing.T) {
	assert := assert.New(t)
	msg := &message.Msg{Type: 2}
	b, err := protojson.Marshal(msg)
	assert.NoError(err)
	v := string(b)
	fmt.Println("default marshalling with known enum", v)
	assert.Contains(v, `:"TYPE_TWO"`)

	msg = &message.Msg{Type: 3}
	b, err = protojson.Marshal(msg)
	assert.NoError(err)
	v = string(b)
	fmt.Println("default marshalling with unknown enum", v)
	assert.Contains(v, `:3`)
}
