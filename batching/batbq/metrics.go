package batbq

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// metrics labels
const (
	BATBQ   = "batbq"
	BATCHER = "batcher"
)

// Metrics stores Batcher Metrics.
type Metrics struct {
	// State
	NumWorkers           *prometheus.GaugeVec
	PendingMessages      *prometheus.GaugeVec
	PendingConfirmations *prometheus.GaugeVec

	// Results
	ReceivedMessages  *prometheus.CounterVec
	ProcessedMessages *prometheus.CounterVec
	ProcessedBatches  *prometheus.CounterVec
	InsertErrors      *prometheus.CounterVec

	// Latencies
	InsertLatency *prometheus.HistogramVec
	AckLatency    *prometheus.HistogramVec
}

func newMetrics(prefix ...string) *Metrics {
	ns := strings.Join(prefix, "_")
	label := []string{BATCHER}
	return &Metrics{
		// State
		NumWorkers: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name:      "workers",
			Namespace: ns,
			Subsystem: BATBQ,
		}, label),
		PendingMessages: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name:      "pending_messages",
			Namespace: ns,
			Subsystem: BATBQ,
		}, label),
		PendingConfirmations: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name:      "pending_confirmations",
			Namespace: ns,
			Subsystem: BATBQ,
		}, label),

		// Results
		ReceivedMessages: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name:      "received_messages_total",
			Namespace: ns,
			Subsystem: BATBQ,
		}, label),
		ProcessedMessages: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name:      "processed_messages_total",
			Namespace: ns,
			Subsystem: BATBQ,
		}, label),
		ProcessedBatches: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name:      "processed_batches_total",
			Namespace: ns,
			Subsystem: BATBQ,
		}, label),
		InsertErrors: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name:      "insert_errors_total",
			Namespace: ns,
			Subsystem: BATBQ,
		}, label),

		// Latencies
		InsertLatency: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name:      "insert_latency_seconds",
			Namespace: ns,
			Subsystem: BATBQ,
		}, label),
		AckLatency: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name:      "ack_latency_seconds",
			Namespace: ns,
			Subsystem: BATBQ,
		}, label),
	}
}

// Register registers all metrics.
func (m *Metrics) Register(reg prometheus.Registerer) {
	reg.MustRegister(m.NumWorkers)
	reg.MustRegister(m.PendingMessages)
	reg.MustRegister(m.PendingConfirmations)

	reg.MustRegister(m.ReceivedMessages)
	reg.MustRegister(m.ProcessedBatches)
	reg.MustRegister(m.ProcessedMessages)
	reg.MustRegister(m.InsertErrors)

	reg.MustRegister(m.InsertLatency)
	reg.MustRegister(m.AckLatency)
}

// Watch registers the metrics and continuously logs them to the console.
func (m *Metrics) Watch(ctx context.Context) {
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
