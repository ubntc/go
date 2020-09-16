package batbq

import (
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

// metrics labels
const (
	BatBq    = "batbq"
	Pipeline = "pipeline"
)

// Metrics stores Batcher Metrics.
type Metrics struct {
	NumWorkers        *prometheus.GaugeVec
	ReceivedMessages  *prometheus.CounterVec
	ProcessedMessages *prometheus.CounterVec
	ProcessedBatches  *prometheus.CounterVec
	InsertErrors      *prometheus.CounterVec
	InsertLatency     *prometheus.HistogramVec
	AckLatency        *prometheus.HistogramVec
}

func newMetrics(prefix ...string) *Metrics {
	ns := strings.Join(prefix, "_")
	label := []string{Pipeline}
	return &Metrics{
		NumWorkers: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name:      "workers",
			Namespace: ns,
			Subsystem: BatBq,
		}, label),
		ProcessedMessages: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name:      "processed_messages_total",
			Namespace: ns,
			Subsystem: BatBq,
		}, label),
		ProcessedBatches: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name:      "processed_batches_total",
			Namespace: ns,
			Subsystem: BatBq,
		}, label),
		ReceivedMessages: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name:      "received_messages_total",
			Namespace: ns,
			Subsystem: BatBq,
		}, label),
		InsertErrors: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name:      "insert_errors_total",
			Namespace: ns,
			Subsystem: BatBq,
		}, label),
		InsertLatency: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name:      "insert_latency_seconds",
			Namespace: ns,
			Subsystem: BatBq,
		}, label),
		AckLatency: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name:      "ack_latency_seconds",
			Namespace: ns,
			Subsystem: BatBq,
		}, label),
	}
}

// Register registers all metrics.
func (m *Metrics) Register(reg prometheus.Registerer) {
	reg.MustRegister(m.NumWorkers)
	reg.MustRegister(m.InsertErrors)
	reg.MustRegister(m.InsertLatency)
	reg.MustRegister(m.AckLatency)
	reg.MustRegister(m.ReceivedMessages)
	reg.MustRegister(m.ProcessedBatches)
	reg.MustRegister(m.ProcessedMessages)
}
