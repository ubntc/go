package batbq_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"github.com/ubntc/go/batching/batbq"

	dummy "github.com/ubntc/go/batching/batbq/_examples/simple/dummy"
	"github.com/ubntc/go/batching/batbq/config"
)

type timing struct {
	dur        time.Duration // batch interval
	sendDelay  time.Duration // delay between test messages
	writeDelay time.Duration // delay of the batch receiver

	autoScale     bool
	scaleInterval time.Duration // how often to trigger worker scaling
}

type testSpec struct {
	len        int   // number of test messages
	cap        int   // batch capacity
	expBatches int   // number of resulting batches
	expErr     error //
	timing     *timing
}

type testResults struct {
	batbq.Metrics
	tableSize int
}

func testRun(t *testing.T, spec testSpec) *testResults {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if spec.timing == nil {
		spec.timing = &timing{time.Second, 0, 0, false, 0}
	}

	data := make([]dummy.Message, 0, spec.len)
	for i := 0; i < spec.len; i++ {
		data = append(data, dummy.Message{ID: fmt.Sprint(i)})
	}
	assert.Len(t, data, spec.len)

	src := &dummy.Source{
		Messages:  data,
		SendDelay: spec.timing.sendDelay,
	}
	input := src.Chan()
	output := &dummy.Putter{
		Name:       t.Name(),
		WriteDelay: spec.timing.writeDelay,
	}

	batcher := batbq.NewInsertBatcher("test", &batbq.WithConfig{
		Capacity:      spec.cap,
		FlushInterval: spec.timing.dur,
		WorkerConfig: config.WorkerConfig{
			ScaleInterval: spec.timing.scaleInterval,
			AutoScale:     spec.timing.autoScale,
		},
	})
	batcher.Process(ctx, input, output)

	return &testResults{
		Metrics:   *batcher.Metrics(),
		tableSize: output.GetLength(),
	}
}

func TestInsertBatcher(t *testing.T) {
	cases := map[string]testSpec{
		// common cases
		"small len": {5, 2, 3, nil, nil},
		"big len":   {1000, 10, 100, nil, nil},
		"small cap": {100, 1, 100, nil, nil},
		"big cap":   {100, 1000, 1, nil, nil},
		// special cases
		"timeout":  {2, 10, 2, nil, &timing{time.Microsecond, time.Millisecond, 0, false, 0}},
		"zero len": {0, 10, 0, nil, nil},
		"zero cap": {10, 0, 10, nil, nil},
		// NOTE: A zero capacity case is valid since Go's `append` will add missing slice capacity
		//       and the feeding of the slice will be done using a zero capacity (blocking) channel.
	}

	for name := range cases {
		t.Run(name, func(t *testing.T) {
			spec := cases[t.Name()]
			res := testRun(t, spec)
			assert.Equal(t, spec.len, res.tableSize)
			// assert.Equal(t, spec.expBatches, ... )         TODO: How to read prometheus metric?
			// assert.GreaterOrEqual(t, res.MaxWorkers, ... ) TODO: How to read prometheus metric?
		})
	}
}

func TestWorkerScaling(t *testing.T) {
	// TODO: add batchDelay to simulate stalled writes
	spec := testSpec{
		100, 10, 10, nil,
		&timing{100 * time.Millisecond, 0, 10 * time.Millisecond, true, time.Millisecond},
	}
	_ = testRun(t, spec)
	// TODO: read prometheus metrics
	// assert.Equal(t, spec.expBatches, int(res.NumBatches))
	// assert.GreaterOrEqual(t, res.MaxWorkers, 2, "at least one extra worker must have started")
}

var testConfig = batbq.WithConfig{Capacity: 10, FlushInterval: 10 * time.Millisecond}

func TestHandleInsertErrors(t *testing.T) {
	ins := batbq.NewInsertBatcher("test", testConfig)
	src := dummy.Source{
		Messages: []dummy.Message{
			{ID: "err1"},
			{ID: "err2"},
			{ID: "ok1"},
		},
	}
	output := &dummy.Putter{}
	ins.Process(context.Background(), src.Chan(), output)
	assert.Len(t, output.InsertErrors, 1)
	assert.Len(t, output.InsertErrors[0], 2)
}

func TestHandleError(t *testing.T) {
	ins := batbq.NewInsertBatcher("test", testConfig)
	src := dummy.Source{
		Messages: []dummy.Message{
			{ID: "fatal"},
			{ID: "ok1"},
			{ID: "ok2"},
		},
	}
	p := &dummy.Putter{}
	ins.Process(context.Background(), src.Chan(), p)
	assert.Len(t, p.InsertErrors, 0)
	assert.Error(t, p.FatalErr, 0)
	assert.Len(t, p.Table, 0)
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
	cfg := config.BatcherConfig{}
	def := cfg.WithDefaults()
	assert.Equal(t, config.BatcherConfig{}, cfg, "orig config must not be modified")
	assert.Equal(t, config.DefaultFlushInterval, def.FlushInterval)
	assert.Equal(t, 1, def.MaxWorkers)
	assert.Equal(t, 1, def.MinWorkers)

	cfg = config.BatcherConfig{}
	cfg.AutoScale = true
	def = cfg.WithDefaults()
	assert.Equal(t, config.DefaultMaxWorkers, def.MaxWorkers)
	assert.Equal(t, config.DefaultMinWorkers, def.MinWorkers)

}
