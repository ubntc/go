package diff_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ubntc/go/playground/diff"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestFunc(t *testing.T) {

	var (
		tsNon = &timestamppb.Timestamp{Seconds: 0, Nanos: 0}
		tsNew = &timestamppb.Timestamp{Seconds: 1, Nanos: 1}
		tsUpd = &timestamppb.Timestamp{Seconds: 2, Nanos: 2}
		tsNil *timestamppb.Timestamp
	)

	changesNew := diff.Diff(tsNon, tsNew)
	changesUpd := diff.Diff(tsNew, tsUpd)
	changesDel := diff.Diff(tsNew, tsNil)

	assert.Equal(t, "created", changesNew["google.protobuf.Timestamp.nanos"])
	assert.Equal(t, "created", changesNew["google.protobuf.Timestamp.seconds"])

	assert.Equal(t, "updated", changesUpd["google.protobuf.Timestamp.nanos"])
	assert.Equal(t, "updated", changesUpd["google.protobuf.Timestamp.seconds"])

	assert.Equal(t, "deleted", changesDel["google.protobuf.Timestamp.nanos"])
	assert.Equal(t, "deleted", changesDel["google.protobuf.Timestamp.seconds"])
}
