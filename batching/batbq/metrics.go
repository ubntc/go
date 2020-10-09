package batbq

import (
	"strings"

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
