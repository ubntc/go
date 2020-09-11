package batbq_test

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"cloud.google.com/go/bigquery"
	"github.com/stretchr/testify/assert"
	"github.com/ubntc/go/batching/batbq"
)

type row struct {
	ID string
}

func (m *row) Save() (row map[string]bigquery.Value, insertID string, err error) {
	v := bigquery.Value(m.ID)
	return map[string]bigquery.Value{"id": v}, m.ID, nil
}

type putter struct {
	sync.Mutex
	table      []map[string]bigquery.Value
	writeDelay time.Duration
	numBatches int
	maxWorkers int
}

func (p *putter) Put(ctx context.Context, src interface{}) error {
	rows := src.([]*bigquery.StructSaver)
	time.Sleep(p.writeDelay)
	p.Lock()
	defer p.Unlock()
	for _, v := range rows {
		row, _, _ := v.Save()
		p.table = append(p.table, row)
	}
	p.numBatches++
	return nil
}

func (p *putter) NumBatches() int {
	p.Lock()
	defer p.Unlock()
	return p.numBatches
}

func (p *putter) Length() int {
	p.Lock()
	defer p.Unlock()
	return len(p.table)
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
		// log.Print("send chan closed")
	}()
	return ch
}

type Timing struct {
	dur           time.Duration // batch interval
	sendDelay     time.Duration // delay between test messages
	writeDelay    time.Duration // delay of the batch receiver
	scaleInterval time.Duration // how often to trigger worker scaling
}

type TestSpec struct {
	len        int   // number of test messages
	cap        int   // batch capacity
	expBatches int   // number of resulting batches
	expErr     error //
	timing     *Timing
}

type TestResults struct {
	batbq.Metrics
	tableSize int
}

func testRun(t *testing.T, spec TestSpec) *TestResults {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if spec.timing == nil {
		spec.timing = &Timing{time.Second, 0, 0, 0}
	}

	data := make([]*batbq.LogMessage, 0, spec.len)
	for i := 0; i < spec.len; i++ {
		row := bigquery.StructSaver{InsertID: fmt.Sprint(i)}
		data = append(data, &batbq.LogMessage{&row})
	}
	assert.Len(t, data, spec.len)

	snd := &source{
		messages:  data,
		sendDelay: spec.timing.sendDelay,
	}
	input := snd.Chan()
	output := &putter{
		writeDelay: spec.timing.writeDelay,
	}

	batcher := batbq.NewInsertBatcher(batbq.BatcherConfig{
		Capacity:      spec.cap,
		FlushInterval: spec.timing.dur,
		ScaleInterval: spec.timing.scaleInterval,
	})
	batcher.Process(ctx, input, output)

	return &TestResults{
		Metrics:   *batcher.Metrics(),
		tableSize: output.Length(),
	}
}

func TestInsertBatcher(t *testing.T) {
	cases := map[string]TestSpec{
		// common cases
		"small len": {5, 2, 3, nil, nil},
		"big len":   {1000, 10, 100, nil, nil},
		"small cap": {100, 1, 100, nil, nil},
		"big cap":   {100, 1000, 1, nil, nil},
		// special cases
		"timeout":  {2, 10, 2, nil, &Timing{time.Microsecond, time.Millisecond, 0, 0}},
		"zero len": {0, 10, 0, nil, nil},
		"zero cap": {10, 0, 10, nil, nil},
		// NOTE: A zero capacity case is valid since Go's `append` will add missing slice capacity
		//       and the feeding of the slice will be done using a zero capacity (blocking) channel.
	}

	for name := range cases {
		t.Run(name, func(t *testing.T) {
			spec := cases[name]
			res := testRun(t, spec)
			assert.Equal(t, spec.len, res.tableSize)
			assert.Equal(t, spec.expBatches, int(res.NumBatches))
			assert.GreaterOrEqual(t, res.MaxWorkers, 1)
		})
	}
}

func TestWorkerScaling(t *testing.T) {
	// TODO: add batchDelay to simulate stalled writes
	spec := TestSpec{100, 10, 10, nil, &Timing{100 * time.Millisecond, 0, 10 * time.Millisecond, time.Millisecond}}
	res := testRun(t, spec)
	assert.Equal(t, spec.expBatches, int(res.NumBatches))
	assert.GreaterOrEqual(t, res.MaxWorkers, 2, "at least one extra worker must have started")
}
