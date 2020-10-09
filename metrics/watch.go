package metrics

import (
	"context"
	"log"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// Metrics represents a set of registerable metrics.
type Metrics interface {
	Register(reg prometheus.Registerer)
}

// Watch registers the metrics and continuously logs them to the console.
func Watch(ctx context.Context, m Metrics) {
	go func() {
		tick := time.NewTicker(time.Second)
		defer tick.Stop()
		log.Print("start watching metrics")
		defer log.Print("stopped watching metrics")
		reg := prometheus.NewPedanticRegistry()
		m.Register(reg)
		for {
			select {
			case <-ctx.Done():
				return
			case <-tick.C:
				printMetrics(reg)
			}
		}
	}()
}

// printMetrics gathers and prints registered metrics.
func printMetrics(reg prometheus.Gatherer) {
	metrics, _ := reg.Gather()
	for _, mf := range metrics {
		for _, m := range mf.GetMetric() {
			val := 0.0
			count := m.GetCounter().GetValue()
			gauge := m.GetGauge().GetValue()
			buckets := m.GetHistogram().GetBucket()

			switch {
			case gauge != 0:
				val = gauge
			case count != 0:
				val = count
			case len(buckets) > 0:
				val = buckets[0].GetExemplar().GetValue()
			}

			if val != 0 {
				labels := m.GetLabel()
				label := ""
				if len(labels) > 0 {
					label = labels[0].GetValue()
				}
				log.Printf("%s[%s] %f", mf.GetName(), label, val)
			}
		}
	}
}
