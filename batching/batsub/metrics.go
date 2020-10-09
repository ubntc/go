package batsub

import (
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

// metrics labels
const (
	BATPS        = "batps"
	SUBSCRIPTION = "subscription"
)

// Metrics stores Batcher Metrics.
type Metrics struct {
	// State
	PendingMessages *prometheus.GaugeVec

	// Results
	ProcessedMessages *prometheus.CounterVec
	ProcessedBatches  *prometheus.CounterVec

	// Latencies
	ProcessingLatency *prometheus.HistogramVec
}

func newMetrics(prefix ...string) *Metrics {
	ns := strings.Join(prefix, "_")
	label := []string{SUBSCRIPTION}
	return &Metrics{
		// State
		PendingMessages: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name:      "pending_messages",
			Namespace: ns,
			Subsystem: BATPS,
		}, label),
		// Results
		ProcessedMessages: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name:      "processed_messages_total",
			Namespace: ns,
			Subsystem: BATPS,
		}, label),
		ProcessedBatches: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name:      "processed_batches_total",
			Namespace: ns,
			Subsystem: BATPS,
		}, label),

		// Latencies
		ProcessingLatency: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name:      "processing_latency_seconds",
			Namespace: ns,
			Subsystem: BATPS,
		}, label),
	}
}

// Register registers all metrics.
func (m *Metrics) Register(reg prometheus.Registerer) {
	reg.MustRegister(m.PendingMessages)
	reg.MustRegister(m.ProcessedBatches)
	reg.MustRegister(m.ProcessedMessages)
	reg.MustRegister(m.ProcessingLatency)
}
