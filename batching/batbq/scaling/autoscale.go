package scaling

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/ubntc/go/batching/batbq/config"
)

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
	s.Lock()
	defer s.Unlock()

	outgoing := float64(batchSize)
	incoming := float64(pendingSize)
	cap := float64(capacity)
	// cap80 := cap * 0.8
	cap20 := cap * 0.2

	switch {
	case outgoing < cap:
		// The workers are not able to fill the batches. This can be caused two causes.
		switch {
		case incoming < cap20:
			// 1. There are just not enough incoming messages.
			s.loadLevel--
		default:
			// 2. There are enough incoming messages. But we cannot tell if more workers would
			//    help or harm, since we do not know if we hit the CPU bounds.
			// TODO: Check CPU load here?
		}
	case outgoing == cap:
		// The batches are full, more workers could help.
		s.loadLevel++
	}
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
		obs     = config.DefaultScaleObservations
		dur     = cfg.ScaleInterval
		secs    = dur / time.Second
		tick    = time.NewTicker(dur)
		highObs = obs / 2 // scale up quickly
		lowObs  = obs * 2 // scale down later
	)

	go func() {
		if cfg.Capacity <= 1 {
			log.Print("skipping to start capacity-based autoscaling for capacity <= 1")
			return
		}
		log.Print(
			"starting autoscaler:\n",
			fmt.Sprintf("-> scaling up after %d continuous high load observations in %ds\n", highObs, secs),
			fmt.Sprintf("-> scaling down after %d continuous low load observations in %ds\n", lowObs, secs),
		)
		for {
			<-tick.C
			if status.Get() > highObs {
				addWorker()
				status.Reset()
			}
			if status.Get() < -lowObs {
				rmWorker()
				status.Reset()
			}
		}
	}()

	wg.Wait()   // wait for all workers to finish
	tick.Stop() // stop worker scaling
}
