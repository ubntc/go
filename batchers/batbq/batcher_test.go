package batbq_test

import (
	"context"
	"fmt"
	"log"
	"sync"
	"testing"
	"time"

	"cloud.google.com/go/bigquery"
	"github.com/stretchr/testify/assert"
	"github.com/ubntc/go/batchers/batbq"
)

type row struct {
	ID string
}

func (m *row) Save() (row map[string]bigquery.Value, insertID string, err error) {
	v := bigquery.Value(m.ID)
	return map[string]bigquery.Value{"id": v}, m.ID, nil
}

type putter struct {
	table      []map[string]bigquery.Value
	numBatches int
}

func (p *putter) Put(ctx context.Context, src interface{}) error {
	rows := src.([]*bigquery.StructSaver)
	for _, v := range rows {
		row, _, _ := v.Save()
		p.table = append(p.table, row)
	}
	p.numBatches++
	return nil
}

type source struct {
	messages  []*batbq.LogMessage
	sendDelay time.Duration
}

func (rec *source) Chan() <-chan batbq.Message {
	ch := make(chan batbq.Message, 100)
	go func() {
		for _, m := range rec.messages {
			ch <- m
			time.Sleep(rec.sendDelay)
		}
		close(ch)
		log.Print("send chan closed")
	}()
	return ch
}

func TestInsertBatcher(t *testing.T) {
	type Spec struct {
		len        int           // number of test messages
		cap        int           // batch capacity
		dur        time.Duration // batch interval
		sendDelay  time.Duration // delay between test messages
		expBatches int           // number of resulting batches
		expErr     error         //
	}
	cases := map[string]Spec{
		// common cases
		"small len": {5, 2, time.Second, 0, 3, nil},
		"big len":   {1000, 10, time.Second, 0, 100, nil},
		"small cap": {100, 1, time.Second, 0, 100, nil},
		"big cap":   {100, 1000, time.Second, 0, 1, nil},
		// special cases
		"timeout":  {2, 10, time.Microsecond, time.Millisecond, 2, nil},
		"zero len": {0, 10, time.Second, 0, 0, nil},
		"zero cap": {10, 0, time.Second, 0, 10, nil},
		// NOTE: A zero capacity case is valid since Go's `append` will add missing slice capacity
		//       and the feeding of the slice will be done using a zero capacity (blocking) channel.
	}

	var wg sync.WaitGroup
	for name, spec := range cases {
		wg.Add(1)
		go func(name string, spec Spec) {
			t.Run(name, func(t *testing.T) {
				defer wg.Done()
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()

				data := make([]*batbq.LogMessage, 0, spec.len)
				for i := 0; i < spec.len; i++ {
					row := bigquery.StructSaver{InsertID: fmt.Sprint(i)}
					data = append(data, &batbq.LogMessage{&row})
				}
				assert.Len(t, data, spec.len)

				snd := &source{
					messages:  data,
					sendDelay: spec.sendDelay,
				}
				input := snd.Chan()
				output := &putter{}
				batcher := batbq.NewInsertBatcher(spec.cap, spec.dur, 1)
				batcher.Process(ctx, input, output)

				assert.Len(t, output.table, spec.len)
				assert.Equal(t, spec.expBatches, output.numBatches)
			})
		}(name, spec)
	}

	wg.Wait()

}
