package config

import "time"

// BatcherConfig defaults.
const (
	DefaultScaleInterval = 3 * time.Second // how often to trigger worker scaling
	DefaultFlushInterval = time.Second     // when to send partially filled batches
	DefaultMinWorkers    = 1
	DefaultMaxWorkers    = 10
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

// WithDefaults copies the config by value, sets missing defaults values returns the copy.
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
	if !cfg.AutoScale {
		cfg.MaxWorkers = 1
		cfg.MinWorkers = 1
	}
	return cfg
}

// Default returns a new default config.
func Default() BatcherConfig {
	cfg := BatcherConfig{}
	return cfg.WithDefaults()
}
