package batbq

import "github.com/ubntc/go/batching/batbq/config"

// BatcherOption configures the batcher.
type BatcherOption interface {
	apply(*InsertBatcher)
}

// Config wraps a config.BatcherConfig to be used as BatcherOption.
type Config config.BatcherConfig

// apply applies the config and sets defaults.
func (cfg Config) apply(ins *InsertBatcher) {
	ins.cfg = config.BatcherConfig(cfg).WithDefaults()
}

func (m *Metrics) apply(ins *InsertBatcher) {
	ins.metrics = m
}
