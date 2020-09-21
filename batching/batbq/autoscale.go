package batbq

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

// scalingStatus safely tracks the load level and scaling status.
type scalingStatus struct {
	loadLevel int
	sync.Mutex
}

func (s *scalingStatus) Inc() {
	s.Lock()
	defer s.Unlock()
	s.loadLevel++
}

func (s *scalingStatus) Dec() {
	s.Lock()
	defer s.Unlock()
	s.loadLevel--
}

func (s *scalingStatus) Reset() {
	s.Lock()
	defer s.Unlock()
	s.loadLevel = 0
}

func (s *scalingStatus) Get() int {
	s.Lock()
	defer s.Unlock()
	return s.loadLevel
}

// autoscale start and stops workers according to number of `ins.cfg.MinWorkers`,
// `ins.cfg.MaxWorkerFactor`, and the number of queued messages on the `input` channel.
func (ins *InsertBatcher) autoscale(ctx context.Context) {
	var wg sync.WaitGroup

	cfg := ins.cfg
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

			ins.worker(wctx, workerNum)

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
		obs     = DefaultScaleObservations
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
			if ins.scaling.Get() > highObs {
				addWorker()
				ins.scaling.Reset()
			}
			if ins.scaling.Get() < -lowObs {
				rmWorker()
				ins.scaling.Reset()
			}
		}
	}()

	wg.Wait()   // wait for all workers to finish
	tick.Stop() // stop worker scaling
}
