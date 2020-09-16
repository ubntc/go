package batbq

import "time"

// BatcherConfig defaults.
const (
	DefaultScaleInterval = time.Second // how often to trigger worker scaling
	DefaultFlushInterval = time.Second // when to send partially filled batches
	DefaultMinWorkers    = 1
	DefaultMaxWorkers    = 10

	DrainedDivisor = 10 // divisor applied to input channel length to check for drained channels
)

// WorkerConfig defines how many workers to use.
type WorkerConfig struct {
	MinWorkers    int
	MaxWorkers    int
	AutoScale     bool
	ScaleInterval time.Duration
}

// BatcherConfig stores InsertBatcher paramaters.
type BatcherConfig struct {
	Capacity      int
	FlushInterval time.Duration
	WorkerConfig
}

// Apply sets the batchers config.
func (cfg BatcherConfig) Apply(ins *InsertBatcher) {
	ins.cfg = cfg.WithDefaults()
}

// WithDefaults loads defaults values for unset values and returns the merged config.
func (cfg BatcherConfig) WithDefaults() BatcherConfig {
	if cfg.FlushInterval <= 0 {
		cfg.FlushInterval = DefaultFlushInterval
	}
	if cfg.MinWorkers <= 0 {
		cfg.MinWorkers = DefaultMinWorkers
	}
	if cfg.MaxWorkers <= 0 {
		cfg.MaxWorkers = DefaultMaxWorkers
	}
	if cfg.ScaleInterval <= 0 {
		cfg.ScaleInterval = DefaultScaleInterval
	}
	return cfg
}

// WithMetrics sets the batchers metrics.
type WithMetrics struct {
	*Metrics
}

// Apply applies the option.
func (m *WithMetrics) Apply(ins *InsertBatcher) {
	ins.metrics = m.Metrics
}
