package batbq_test

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"cloud.google.com/go/bigquery"
	"github.com/prometheus/client_golang/prometheus"
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
	name string
	sync.Mutex
	table        []map[string]bigquery.Value
	writeDelay   time.Duration
	numBatches   int
	maxWorkers   int
	numErrors    int
	insertErrors []error
	fatalErr     error
}

func (p *putter) Put(ctx context.Context, src interface{}) error {
	rows := src.([]*bigquery.StructSaver)
	time.Sleep(p.writeDelay)
	p.Lock()
	defer p.Unlock()
	errors := make(bigquery.PutMultiError, 0)
	for i, v := range rows {
		row, insertID, err := v.Save()
		if insertID == "fatal" {
			p.fatalErr = fmt.Errorf("all inserts failed")
			return p.fatalErr
		}
		if len(insertID) >= 3 && insertID[:3] == "err" || err != nil {
			errors = append(errors, bigquery.RowInsertionError{RowIndex: i, InsertID: insertID})
			continue
		}
		p.table = append(p.table, row)
	}
	p.numBatches++
	if len(errors) > 0 {
		p.insertErrors = append(p.insertErrors, errors)
		return errors
	}
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
	dur        time.Duration // batch interval
	sendDelay  time.Duration // delay between test messages
	writeDelay time.Duration // delay of the batch receiver

	autoScale     bool
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

// label used to read prometheus metrics
var testLabel = prometheus.Labels{batbq.Pipeline: "test"}

func testRun(t *testing.T, spec TestSpec) *TestResults {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if spec.timing == nil {
		spec.timing = &Timing{time.Second, 0, 0, false, 0}
	}

	data := make([]*batbq.LogMessage, 0, spec.len)
	for i := 0; i < spec.len; i++ {
		row := bigquery.StructSaver{InsertID: fmt.Sprint(i)}
		data = append(data, &batbq.LogMessage{row})
	}
	assert.Len(t, data, spec.len)

	src := &source{
		messages:  data,
		sendDelay: spec.timing.sendDelay,
	}
	input := src.Chan()
	output := &putter{
		name:       t.Name(),
		writeDelay: spec.timing.writeDelay,
	}

	batcher := batbq.NewInsertBatcher("test", batbq.BatcherConfig{
		Capacity:      spec.cap,
		FlushInterval: spec.timing.dur,
		WorkerConfig: batbq.WorkerConfig{
			ScaleInterval: spec.timing.scaleInterval,
			AutoScale:     spec.timing.autoScale,
		},
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
		"timeout":  {2, 10, 2, nil, &Timing{time.Microsecond, time.Millisecond, 0, false, 0}},
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
			// assert.Equal(t, spec.expBatches, ... )         TODO: How to read prometheus metric?
			// assert.GreaterOrEqual(t, res.MaxWorkers, ... ) TODO: How to read prometheus metric?
		})
	}
}

func TestWorkerScaling(t *testing.T) {
	// TODO: add batchDelay to simulate stalled writes
	spec := TestSpec{
		100, 10, 10, nil,
		&Timing{100 * time.Millisecond, 0, 10 * time.Millisecond, true, time.Millisecond},
	}
	_ = testRun(t, spec)
	// TODO: read prometheus metrics
	// assert.Equal(t, spec.expBatches, int(res.NumBatches))
	// assert.GreaterOrEqual(t, res.MaxWorkers, 2, "at least one extra worker must have started")
}

var testConfig = batbq.BatcherConfig{Capacity: 10, FlushInterval: 10 * time.Millisecond}

func TestHandleInsertErrors(t *testing.T) {
	ins := batbq.NewInsertBatcher("test", testConfig)
	src := source{
		messages: []*batbq.LogMessage{
			{bigquery.StructSaver{InsertID: "err1"}},
			{bigquery.StructSaver{InsertID: "err2"}},
			{bigquery.StructSaver{InsertID: "ok1"}},
		},
	}
	output := &putter{}
	ins.Process(context.Background(), src.Chan(), output)
	assert.Len(t, output.insertErrors, 1)
	assert.Len(t, output.insertErrors[0], 2)
}

func TestHandleError(t *testing.T) {
	ins := batbq.NewInsertBatcher("test", testConfig)
	src := source{
		messages: []*batbq.LogMessage{
			{bigquery.StructSaver{InsertID: "fatal"}},
			{bigquery.StructSaver{InsertID: "ok1"}},
			{bigquery.StructSaver{InsertID: "ok2"}},
		},
	}
	p := &putter{}
	ins.Process(context.Background(), src.Chan(), p)
	assert.Len(t, p.insertErrors, 0)
	assert.Error(t, p.fatalErr, 0)
	assert.Len(t, p.table, 0)
}

func TestBatcherWithMetrics(t *testing.T) {
	ins := batbq.NewInsertBatcher("test", &batbq.WithMetrics{})
	assert.NotNil(t, ins.Metrics())
}

func TestBatcherRegisterMetrics(t *testing.T) {
	mtx := batbq.NewInsertBatcher("test").Metrics()
	reg := prometheus.NewPedanticRegistry()
	mtx.Register(reg)
}

func TestDefaults(t *testing.T) {
	cfg := batbq.BatcherConfig{}
	def := cfg.WithDefaults()
	assert.Equal(t, batbq.BatcherConfig{}, cfg, "orig config must be modified")
	assert.Equal(t, batbq.DefaultMinWorkers, def.MinWorkers)
	assert.Equal(t, batbq.DefaultMaxWorkers, def.MaxWorkers)
	assert.Equal(t, batbq.DefaultFlushInterval, def.FlushInterval)
}
