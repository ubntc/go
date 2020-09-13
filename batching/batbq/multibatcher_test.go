package batbq_test

import (
	"context"
	"testing"
	"time"

	"cloud.google.com/go/bigquery"
	"github.com/stretchr/testify/assert"
	"github.com/ubntc/go/batching/batbq"
	custom "github.com/ubntc/go/batching/batbq/_examples/simple/dummy"
)

func TestMultiBatcher(t *testing.T) {
	mb := batbq.NewMultiBatcher(
		[]string{"p1", "p2"},
		batbq.BatcherConfig{
			AutoScale: false,
		},
	)

	input := func(id batbq.PipelineID) <-chan batbq.Message {
		src := custom.NewSource(string(id))
		in := make(chan batbq.Message, 10)
		go func() {
			defer close(in)
			src.Receive(context.Background(), func(m *custom.Message) {
				in <- &batbq.LogMessage{&bigquery.StructSaver{
					InsertID: "id",
					Struct:   custom.Message{ID: "id", Val: 1},
				}}
			})
		}()
		return in
	}

	putters := make(chan *putter, 100)

	output := func(id batbq.PipelineID) batbq.Putter {
		p := &putter{
			name:       string(id),
			writeDelay: time.Microsecond,
		}
		putters <- p
		return p
	}

	mb.Process(context.Background(), input, output)
	close(putters)

	for p := range putters {
		assert.Equal(t, 200, p.Length())
	}
}