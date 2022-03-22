package metrics

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
)

func Test_printMetrics(t *testing.T) {
	c := prometheus.NewCounter(prometheus.CounterOpts{Name: "count"})
	reg := prometheus.NewPedanticRegistry()
	reg.Register(c)
	c.Inc()
	printMetrics(reg)
}
