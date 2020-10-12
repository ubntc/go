package multibatcher_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ubntc/go/batching/batbq"
	dummy "github.com/ubntc/go/batching/batbq/_examples/simple/dummy"
	mb "github.com/ubntc/go/batching/batbq/multibatcher"
)

func dummyData(topic string, size int) []dummy.Message {
	res := make([]dummy.Message, size)
	for i := 0; i < size; i++ {
		res[i] = dummy.Message{ID: string(topic) + "_msg_" + fmt.Sprint(i), Val: i}
	}
	return res
}

// GetInput demonstrates how to implement an InputGetter.
func GetInput(id string) <-chan batbq.Message {
	src := dummy.NewSource(string(id))
	src.Messages = dummyData(id, 10)
	return src.Chan()
}

// GetInput demonstrates how to implement an OutputGetter.
func GetOutput(id string) batbq.Putter {
	p := &dummy.Putter{Name: string(id)}
	return p
}

func TestMultiBatcher(t *testing.T) {
	testPutters := make(chan *dummy.Putter, 100)

	// Example MultiBatcher Setup
	// ==========================
	// 1. Create a batcher with topic/table names.
	mb := mb.NewMultiBatcher(
		[]string{"p1", "p2"},
	)

	// 2. Implement an InputGetter that returns the input chan `<-batbq.Message`.
	getInput := GetInput

	// 3. Implement an OutputGetter that returns an output `batbq.Putter`.
	getOutput := func(id string) batbq.Putter {
		p := GetOutput(id)
		testPutters <- p.(*dummy.Putter)
		return p
	}

	// 4. Start the multi batcher with mb.Process or mb.MustProcess.
	err := mb.MustProcess(context.Background(), getInput, getOutput)

	assert.NoError(t, err)
	close(testPutters)
	for p := range testPutters {
		assert.Equal(t, 10, p.GetLength())
	}
}
