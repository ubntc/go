package batbq

import (
	"context"
	"log"
	"sync"
	"time"
)

// Putter provides a `Put` func as used by the `bigquery.Inserter`.
type Putter interface {
	Put(ctx context.Context, src interface{}) error
}

// InsertBatcher implements automatic batching with a batch capacity and flushInterval.
type InsertBatcher struct {
	cfg     BatcherConfig
	metrics *metricsRecorder
}

// NewInsertBatcher returns an InsertBatcher.
func NewInsertBatcher(cfg BatcherConfig) *InsertBatcher {
	return &InsertBatcher{cfg.WithDefaults(), newMetricsRecorder()}
}

// Metrics returns a copy of the metrics.
func (ins *InsertBatcher) Metrics() *Metrics {
	return ins.metrics.Metrics()
}

// Process starts the batcher.
func (ins *InsertBatcher) Process(ctx context.Context, input <-chan Message, out Putter) {
	if ins.cfg.AutoScale {
		ins.autoProcess(ctx, input, out)
		return
	}
	ins.metrics.SetWorkers(1)
	ins.process(ctx, input, out)
	ins.metrics.SetWorkers(0)
}

func (ins *InsertBatcher) process(ctx context.Context, input <-chan Message, out Putter) {
	var wg sync.WaitGroup
	defer wg.Wait()

	// batch of messages to be filled from the input channel
	var batch []Message
	cfg := ins.cfg

	var ackLock sync.Mutex
	confirm := func(messages []Message, err error) {
		ackLock.Lock()
		// The lock ensures to stop processing if the previous acks are not complete.
		// Otherwise the batcher could eat up the memory with pending unacked messages
		// in case the (n)acking takes too long.

		// Also ensure we wait for pending (n)acks after the batcher stops.
		wg.Add(1)
		go func() {
			defer wg.Done()        // allow the batcher to stop
			defer ackLock.Unlock() // allow confirming the next batch
			// TODO: handle insert errors
			if err != nil {
				log.Print(err)
			}
			for _, m := range messages {
				m.Ack()
			}
			ins.metrics.IncMessages(len(messages), 1)
		}()
	}

	ticker := time.NewTicker(cfg.FlushInterval)
	defer ticker.Stop()

	flush := func() {
		if len(batch) == 0 {
			return
		}
		err := out.Put(ctx, toStructs(batch))
		confirm(batch, err)                      // use current slice to ack/nack processed messages
		batch = make([]Message, 0, cfg.Capacity) // create a new slice to allow immediate refill
		ticker.Reset(cfg.FlushInterval)          // reset the ticker to avoid too early flush
	}
	defer flush()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			flush()
		case msg, more := <-input:
			if !more {
				return
			}
			batch = append(batch, msg)
			if len(batch) >= cfg.Capacity {
				flush()
			}
		}
	}
}

// autoProcess start and stops `process` multiple times according to number of `ins.cfg.MinWorkers`,
// `ins.cfg.MaxWorkerFactor`, and the number of queued messages on the `input` channel.
func (ins *InsertBatcher) autoProcess(ctx context.Context, input <-chan Message, out Putter) {
	var wg sync.WaitGroup

	cfg := ins.cfg
	hooks := make(map[context.Context]func())
	mu := &sync.Mutex{}

	addWorker := func() {
		mu.Lock()
		defer mu.Unlock()

		if len(hooks) >= cfg.MinWorkers*MaxWorkerFactor {
			return
		}
		log.Printf("adding worker #%d", len(hooks)+1)
		wctx, cancel := context.WithCancel(ctx)
		hooks[wctx] = cancel
		ins.metrics.SetWorkers(len(hooks))

		wg.Add(1)
		go func() {
			defer wg.Done()

			ins.process(wctx, input, out)

			mu.Lock()
			delete(hooks, wctx)
			ins.metrics.SetWorkers(len(hooks))
			mu.Unlock()
		}()
	}

	rmWorker := func() {
		mu.Lock()
		defer mu.Unlock()
		if len(hooks) <= cfg.MinWorkers {
			return
		}
		log.Printf("removing one of %d workers", len(hooks))
		for _, cancel := range hooks {
			cancel()
			break
		}
	}

	for i := 0; i < cfg.MinWorkers; i++ {
		addWorker()
	}

	// start worker scaling
	tick := time.NewTicker(cfg.ScaleInterval)

	go func() {
		if cfg.Capacity <= 1 {
			// cannot do capacity-based scaling for small capacities
			return
		}
		for {
			<-tick.C
			switch {
			case len(input) >= cfg.Capacity:
				addWorker()
			case len(input) < cfg.Capacity/DrainedDivisor:
				rmWorker()
			}
		}
	}()

	wg.Wait()   // wait for all workers to finish
	tick.Stop() // stop worker scaling
}
