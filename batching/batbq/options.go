package batbq

import "github.com/ubntc/go/batching/batbq/config"

// BatcherOption configures the batcher.
type BatcherOption interface {
	Apply(*InsertBatcher)
}

// WithConfig set config values.
type WithConfig config.BatcherConfig

// Apply applies the config and sets defaults.
func (cfg WithConfig) Apply(ins *InsertBatcher) {
	ins.cfg = config.BatcherConfig(cfg).WithDefaults()
}

// WithMetrics sets the batchers metrics.
type WithMetrics struct {
	*Metrics
}

// Apply applies the option.
func (m *WithMetrics) Apply(ins *InsertBatcher) {
	ins.metrics = m.Metrics
}
