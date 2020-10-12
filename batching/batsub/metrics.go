package batsub

import (
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

// metrics prefix and label
const (
	BATSUB       = "batsub"
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

// NewMetrics returns prefixed metrics.
func NewMetrics(prefix ...string) *Metrics {
	ns := strings.Join(prefix, "_")
	if len(ns) == 0 {
		ns = BATSUB
	}
	label := []string{SUBSCRIPTION}
	return &Metrics{
		// State
		PendingMessages: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name:      "pending_messages",
			Namespace: ns,
		}, label),

		// Results
		ProcessedMessages: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name:      "processed_messages_total",
			Namespace: ns,
		}, label),
		ProcessedBatches: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name:      "processed_batches_total",
			Namespace: ns,
		}, label),

		// Latencies
		ProcessingLatency: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name:      "processing_latency_seconds",
			Namespace: ns,
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
