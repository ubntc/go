package cloudtables

import (
	"encoding/json"
	"errors"
	"os"
	"strconv"
	"time"
)

const (
	DefaultNumPartitions      = 10
	DefaultReplicationFactor  = 3
	DefaultSegmentBytes       = 50 * 1024 * 1024 // 50MB (min allowed by Confluent Cloud)
	DefaultMaxCompactionLagMs = 7 * 24 * time.Hour / time.Millisecond
	DefaultDeleteRetentionMs  = 24 * time.Hour / time.Millisecond

	DefaultPrefix       = "tables."
	DefaultManagerTopic = DefaultPrefix + "tables"
)

func DefaultCompactConfig() map[string]string {
	segmentBytes := strconv.Itoa(DefaultSegmentBytes)
	maxLagMs := strconv.Itoa(int(DefaultMaxCompactionLagMs))
	deleteMs := strconv.Itoa(int(DefaultDeleteRetentionMs))
	return map[string]string{
		"retention.ms":          "-1",         // Keep data indefinitely
		"cleanup.policy":        "compact",    // and make it a table.
		"min.compaction.lag.ms": "0",          // Also enable table-like instant updates
		"max.compaction.lag.ms": maxLagMs,     // and enforce compaction after 7 days.
		"segment.bytes":         segmentBytes, // Use normal segments size to avoid compaction pressure.
		"delete.retention.ms":   deleteMs,     // Keep tombstones for 24 hours to inform about deletion.
	}
}

type KeyFile struct {
	Key    string `json:"key,omitempty"`
	Secret string `json:"secret,omitempty"`
	Server string `json:"server,omitempty"`
}

func (kf *KeyFile) Load(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, kf)
}

func (kf *KeyFile) Validate() error {
	if len(kf.Key) == 0 {
		return errors.New("key is empty")
	}
	if len(kf.Secret) == 0 {
		return errors.New("secret is empty")
	}
	if len(kf.Server) == 0 {
		return errors.New("server is empty")
	}
	return nil
}

func LoadConfig() (*KeyFile, error) {
	cfg := KeyFile{}

	if err := cfg.Load(os.Getenv("CC_API_KEY_FILE")); err != nil {
		return nil, err
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}
