package batbq

import "time"

// BatcherConfig defaults.
const (
	DefaultScaleInterval = time.Second // how often to trigger worker scaling
	DefaultFlushInterval = time.Second // when to send partially filled batches
	DefaultMinWorkers    = 1

	MaxWorkerFactor = 10 // factor multiplied to cfg.MinWorkers to determine the MaxWorkers
	DrainedDivisor  = 10 // divisor applied to input channel length to check for drained channels
)

// BatcherConfig stores InsertBatcher paramaters.
type BatcherConfig struct {
	Capacity      int
	FlushInterval time.Duration
	MinWorkers    int
	AutoScale     bool
	ScaleInterval time.Duration
}

// WithDefaults loads defaults values for unset values and returns the merged config.
func (cfg BatcherConfig) WithDefaults() BatcherConfig {
	if cfg.FlushInterval <= 0 {
		cfg.FlushInterval = DefaultFlushInterval
	}
	if cfg.MinWorkers <= 0 {
		cfg.MinWorkers = DefaultMinWorkers
	}
	if cfg.ScaleInterval <= 0 {
		cfg.ScaleInterval = DefaultScaleInterval
	}
	return cfg
}
