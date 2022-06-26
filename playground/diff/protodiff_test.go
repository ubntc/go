package diff_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/ubntc/go/playground/diff"
	"github.com/ubntc/go/playground/diff/message"
)

type Map map[string]string

func TestFunc(t *testing.T) {

	var (
		null        *message.Msg
		empty       = &message.Msg{}
		valEmpty    = &message.Msg{Value: ""}
		subNil      = &message.Msg{Sub: nil}
		subEmpty    = &message.Msg{Sub: &message.Sub{}}
		subValEmpty = &message.Msg{Sub: &message.Sub{Value: ""}}

		valX    = &message.Msg{Value: "x"}
		subValX = &message.Msg{Sub: &message.Sub{Value: "x"}}

		valY    = &message.Msg{Value: "y"}
		subValY = &message.Msg{Sub: &message.Sub{Value: "y"}}

		valTs1    = &message.Msg{Ts: &timestamppb.Timestamp{Seconds: 1}}
		subValTs1 = &message.Msg{Sub: &message.Sub{Ts: &timestamppb.Timestamp{Seconds: 1}}}

		valDur = &message.Msg{Dur: durationpb.New(1)}

		valObj = &message.Msg{Obj: &structpb.Struct{Fields: map[string]*structpb.Value{
			"key": {Kind: &structpb.Value_NumberValue{NumberValue: 1.0}},
		}}}
	)

	nilMessages := []*message.Msg{null, empty, valEmpty, subEmpty, subNil, subValEmpty}

	for _, a := range nilMessages {
		for _, b := range nilMessages {
			assert.Nil(t, diff.Diff(a, b))
		}
	}

	assert.Equal(t, Map{"value": "created"}, Map(diff.Diff(null, valX)))
	assert.Equal(t, Map{"sub.value": "created"}, Map(diff.Diff(null, subValX)))
	assert.Equal(t, Map{"ts": "created"}, Map(diff.Diff(null, valTs1)))
	assert.Equal(t, Map{"sub.ts": "created"}, Map(diff.Diff(null, subValTs1)))
	assert.Equal(t, Map{"dur": "created"}, Map(diff.Diff(null, valDur)))
	assert.Equal(t, Map{"obj": "created"}, Map(diff.Diff(null, valObj)))

	assert.Equal(t, Map{"value": "updated"}, Map(diff.Diff(valX, valY)))
	assert.Equal(t, Map{"sub.value": "updated"}, Map(diff.Diff(subValX, subValY)))

	assert.Equal(t, Map{"value": "deleted"}, Map(diff.Diff(valX, null)))
	assert.Equal(t, Map{"value": "deleted"}, Map(diff.Diff(valX, empty)))
	assert.Equal(t, Map{"value": "deleted"}, Map(diff.Diff(valX, valEmpty)))

	assert.Equal(t, Map{"sub.value": "deleted"}, Map(diff.Diff(subValX, empty)))
	assert.Equal(t, Map{"sub.value": "deleted"}, Map(diff.Diff(subValX, subEmpty)))
	assert.Equal(t, Map{"sub.value": "deleted"}, Map(diff.Diff(subValX, subNil)))
	assert.Equal(t, Map{"sub.value": "deleted"}, Map(diff.Diff(subValX, subValEmpty)))
}
