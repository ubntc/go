package batbq

import "sync"

// Metrics stores Batcher Metrics.
type Metrics struct {
	MaxWorkers     int
	CurrentWorkers int
	NumMessages    int64
	NumBatches     int64
}

type metricsRecorder struct {
	m *Metrics
	*sync.Mutex
}

// NewMetricsRecorder returns a new metricsRecorder.
func newMetricsRecorder() *metricsRecorder {
	return &metricsRecorder{
		&Metrics{},
		&sync.Mutex{},
	}
}

// SetWorkers updates the metrics.
func (rec *metricsRecorder) SetWorkers(numWorkers int) {
	rec.Lock()
	defer rec.Unlock()
	rec.m.CurrentWorkers = numWorkers
	if numWorkers > rec.m.MaxWorkers {
		rec.m.MaxWorkers = numWorkers
	}
}

// SetMessages updates the metrics.
func (rec *metricsRecorder) SetMessages(addedMessages, addedBatches int) {
	rec.Lock()
	defer rec.Unlock()
	rec.m.NumBatches += int64(addedBatches)
	rec.m.NumMessages += int64(addedMessages)
}

// Metrics create a copy of the metrics.
func (rec *metricsRecorder) Metrics() *Metrics {
	rec.Lock()
	defer rec.Unlock()
	cp := *rec.m
	return &cp
}
