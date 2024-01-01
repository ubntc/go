package pebble

import (
	"github.com/prometheus/client_golang/prometheus"
	io_prometheus_client "github.com/prometheus/client_model/go"
	"github.com/ubntc/go/kstore/provider/api"
)

type Mtx struct {
	Reads  *prometheus.CounterVec
	Writes *prometheus.CounterVec
}

var Metrics = &Mtx{
	Reads: prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "kstore",
			Subsystem: "pebble",
			Name:      "reads_total",
			Help:      "Messages read by topic and status",
		},
		[]string{"topic", "status"},
	),
	Writes: prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "kstore",
			Subsystem: "pebble",
			Name:      "writes_total",
			Help:      "Messages written by topic",
		},
		[]string{"topic"},
	),
}

func (m *Mtx) ObserveRead(msg api.Message, topic string, status OffsetStatus) {
	m.Reads.WithLabelValues(topic, status.String()).Inc()
}

func (m *Mtx) ObserveWrite(topic string) {
	m.Writes.WithLabelValues(topic).Inc()
}

func (m *Mtx) GetReads(topic string, status ...OffsetStatus) map[OffsetStatus]int {
	if len(status) == 0 {
		status = OffsetStatuses
	}

	result := map[OffsetStatus]int{}

	for _, s := range status {
		result[s] = 0
		metric, err := m.Reads.GetMetricWithLabelValues(topic, s.String())
		if err != nil {
			continue
		}

		pbMetric := &io_prometheus_client.Metric{}
		err = metric.Write(pbMetric)
		if err != nil {
			continue
		}

		result[s] = int(pbMetric.Counter.GetValue())
	}

	// Return the map of offset counts
	return result
}

func (m *Mtx) GetWrites(topic string) int {
	metric, err := m.Writes.GetMetricWithLabelValues(topic)
	if err != nil {
		return 0
	}
	pbMetric := &io_prometheus_client.Metric{}
	if err = metric.Write(pbMetric); err != nil {
		return 0
	}

	return int(pbMetric.Counter.GetValue())
}
