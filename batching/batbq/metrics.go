package batbq

import (
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

// metrics prefix and label name
const (
	BATBQ   = "batbq"
	BATCHER = "batcher"
)

// Metrics stores Batcher Metrics.
type Metrics struct {
	// State
	NumWorkers      *prometheus.GaugeVec
	PendingMessages *prometheus.GaugeVec

	// Results
	ReceivedMessages  *prometheus.CounterVec
	ProcessedMessages *prometheus.CounterVec
	ProcessedBatches  *prometheus.CounterVec
	InsertErrors      *prometheus.CounterVec

	// Latencies
	InsertLatency *prometheus.HistogramVec
	AckLatency    *prometheus.HistogramVec
}

// NewMetrics create returns a new Metrics object.
func NewMetrics(prefix ...string) *Metrics {
	ns := strings.Join(prefix, "_")
	if len(ns) == 0 {
		ns = BATBQ
	}
	label := []string{BATCHER}
	return &Metrics{
		// State
		NumWorkers: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name:      "workers",
			Namespace: ns,
		}, label),
		PendingMessages: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name:      "pending_messages",
			Namespace: ns,
		}, label),

		// Results
		ReceivedMessages: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name:      "received_messages_total",
			Namespace: ns,
		}, label),
		ProcessedMessages: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name:      "processed_messages_total",
			Namespace: ns,
		}, label),
		ProcessedBatches: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name:      "processed_batches_total",
			Namespace: ns,
		}, label),
		InsertErrors: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name:      "insert_errors_total",
			Namespace: ns,
		}, label),

		// Latencies
		InsertLatency: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name:      "insert_latency_seconds",
			Namespace: ns,
		}, label),
		AckLatency: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name:      "ack_latency_seconds",
			Namespace: ns,
		}, label),
	}
}

// Register registers all metrics.
func (m *Metrics) Register(reg prometheus.Registerer) {
	reg.MustRegister(m.NumWorkers)
	reg.MustRegister(m.PendingMessages)

	reg.MustRegister(m.ReceivedMessages)
	reg.MustRegister(m.ProcessedBatches)
	reg.MustRegister(m.ProcessedMessages)
	reg.MustRegister(m.InsertErrors)

	reg.MustRegister(m.InsertLatency)
	reg.MustRegister(m.AckLatency)
}
