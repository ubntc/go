package batsub

import "time"

// Option defines a batching option.
type Option interface {
	apply(sub *BatchedSubscription)
}

func (opt *Metrics) apply(sub *BatchedSubscription) {
	sub.metrics = opt
}

// FlushInterval sets the batcher's flush interval.
type FlushInterval time.Duration

func (opt FlushInterval) apply(sub *BatchedSubscription) {
	sub.flushInterval = time.Duration(opt)
}

// Capacity sets the batcher's capacity.
type Capacity int

func (opt Capacity) apply(sub *BatchedSubscription) {
	sub.capacity = int(opt)
}
