package scaling

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/ubntc/go/batching/batbq/config"
)

// MaxLoadLevel defines the highest possible load level.
const MaxLoadLevel = 10

// MinLoadLevel defines the lowest possible load level.
const MinLoadLevel = -10

// Status safely tracks the load level and scaling status.
type Status struct {
	loadLevel int
	sync.Mutex
}

// Reset resets the load level.
func (s *Status) Reset() {
	s.Lock()
	defer s.Unlock()
	s.loadLevel = 0
}

// Get returns the load level.
func (s *Status) Get() int {
	s.Lock()
	defer s.Unlock()
	return s.loadLevel
}

// UpdateLoadLevel updates the load level based on the observed last batch size compared to the
// configured capacity and the observed number of pending messages compared to the capacity.
//
// The load is considered as high if:
// * the batch size hits the capacity OR
// * the pending size is above 80% of the capacity
//
// The load is considered as low if:
// * the batch size is below the capacity AND
// * the pending size is below 20% of the capacity
//
// The system is considered as overloaded, with more workers being harmful, if:
// * the pending size is above 80% of the capacity AND
// * the batch size is below 80% of the capacity
//
func (s *Status) UpdateLoadLevel(batchSize, pendingSize, capacity int) {
	if s == nil {
		// allow running with empty Status
		return
	}

	s.Lock()
	defer s.Unlock()

	outgoing := float64(batchSize)
	cap := float64(capacity)
	cap50 := cap * 0.5

	// log.Print("update load level", s.loadLevel, batchSize, pendingSize, capacity)

	// Naive Scaling
	// =============
	// 1. Assume that we get enough data and have the CPU to fill the batches up to the capacity.
	// 2. Assume that continuously hitting the limit means that need we need to increase throughput.
	// 3. Assume that more workers will help to fill and send batches concurrently.
	if outgoing >= cap {
		if s.loadLevel == MaxLoadLevel {
			return
		}
		s.loadLevel++
		return
	}

	if outgoing < cap50 {
		if s.loadLevel == MinLoadLevel {
			return
		}
		s.loadLevel--
		return
	}

	/*
		Failed experiment: Determine load from high or low capacity usage

		incoming := float64(pendingSize)
		cap80 := cap * 0.8
		cap50 := cap * 0.5
		cap20 := cap * 0.2

		switch {
		case outgoing < cap80:
			// The workers are not able to fill the batches, which can have two causes:
			switch {
			case incoming < cap20:
				// There are not enough incoming messages.
				s.loadLevel--
			case incoming > cap80:
				// There are too many incoming messages.
				s.loadLevel++ //  TODO: Check CPU load here! (assume overload for now)
			}
		case outgoing > cap80:
			// The batches are very full, which can have two causes:
			switch {
			case incoming > cap80:
				// There are too many incoming messages.
				s.loadLevel++
			case incoming < cap20:
				// There are not enough incoming messages.
				// keep load level, since batches are well filled
			}
		}
	*/
}

// Autoscale starts and stops workers according to the configured `ins.cfg.MinWorkers`,
// `ins.cfg.MaxWorkers`, and the current `ins.scaling.loadLevel`.
// The workers will increase the load level when the batch size hits the capacity and will decrease
// the load level when the batch size is below the capacity.
//
// Autoscaling can be enabled by setting `BatcherConfig.AutoScale = true`.
// See [SCALING.md](../SCALING.md) to check when to use auto scaling.
//
func Autoscale(ctx context.Context, cfg *config.BatcherConfig, status *Status, worker func(ctx context.Context, num int)) {
	var wg sync.WaitGroup

	hooks := make(map[context.Context]func())
	mu := &sync.Mutex{}

	addWorker := func() {
		mu.Lock()
		defer mu.Unlock()

		if len(hooks) >= cfg.MaxWorkers {
			return
		}
		wctx, cancel := context.WithCancel(ctx)
		hooks[wctx] = cancel
		workerNum := len(hooks)

		wg.Add(1)
		go func() {
			defer wg.Done()

			worker(wctx, workerNum)

			mu.Lock()
			delete(hooks, wctx)
			mu.Unlock()
		}()
	}

	rmWorker := func() {
		mu.Lock()
		defer mu.Unlock()
		if len(hooks) <= cfg.MinWorkers {
			return
		}
		log.Printf("removing 1 of %d workers", len(hooks))
		for _, cancel := range hooks {
			cancel()
			break
		}
	}

	for i := 0; i < cfg.MinWorkers; i++ {
		addWorker()
	}

	// start worker scaling
	var (
		dur       = cfg.ScaleInterval
		secs      = dur / time.Second
		tick      = time.NewTicker(dur)
		highLevel = MaxLoadLevel / 2 // scale up quickly
		lowLevel  = MinLoadLevel     // scale down later
	)

	go func() {
		if cfg.Capacity <= 1 {
			log.Print("skipping to start capacity-based autoscaling for capacity <= 1")
			return
		}
		log.Printf("scaling up/down every %ds if load level above %d or below %d",
			secs, highLevel, lowLevel,
		)
		for {
			<-tick.C
			if status.Get() >= highLevel {
				addWorker()
				status.Reset()
			}
			if status.Get() <= lowLevel {
				rmWorker()
				status.Reset()
			}
			// log.Print("load level:", status.Get(), highLevel, lowLevel)
		}
	}()

	wg.Wait()   // wait for all workers to finish
	tick.Stop() // stop worker scaling
}
