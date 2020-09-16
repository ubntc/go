package batbq

import (
	"context"
	"log"
	"sync"
	"time"
)

// autoscale start and stops workers according to number of `ins.cfg.MinWorkers`,
// `ins.cfg.MaxWorkerFactor`, and the number of queued messages on the `input` channel.
func (ins *InsertBatcher) autoscale(ctx context.Context) {
	var wg sync.WaitGroup

	cfg := ins.cfg
	hooks := make(map[context.Context]func())
	mu := &sync.Mutex{}
	input := ins.input

	addWorker := func() {
		mu.Lock()
		defer mu.Unlock()

		if len(hooks) >= cfg.MaxWorkers {
			return
		}
		log.Printf("adding worker #%d", len(hooks)+1)
		wctx, cancel := context.WithCancel(ctx)
		hooks[wctx] = cancel

		wg.Add(1)
		go func() {
			defer wg.Done()

			ins.worker(wctx)

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
