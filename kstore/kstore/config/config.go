package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"strconv"
	"time"
)

var (
	ErrorParseKeyFile = errors.New("Failed to load or parse keyfile")
	ErrorFindKeyFile  = errors.New("Failed to find keyfile")
)

const (
	DefaultNumPartitions      = 10
	DefaultReplicationFactor  = 3
	DefaultSegmentBytes       = 50 * 1024 * 1024 // 50MB (min allowed by Confluent Cloud)
	DefaultMaxCompactionLagMs = 7 * 24 * time.Hour / time.Millisecond
	DefaultDeleteRetentionMs  = 24 * time.Hour / time.Millisecond

	DefaultTopicPrefix  = "tables."
	DefaultSchemasTopic = DefaultTopicPrefix + "schemas"
)

type KafkaProperties map[string]string

func DefaultProperties() map[string]string {
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

type Group struct {
	ID     string
	Topics []string
}

type KeyFile struct {
	Key    string `json:"key,omitempty"`
	Secret string `json:"secret,omitempty"`
	Server string `json:"server,omitempty"`
}

func (kf *KeyFile) Load() error {
	if kf == nil {
		*kf = KeyFile{}
	}
	filename, err := FindKeyFile()
	if err != nil {
		return errors.Join(err, ErrorFindKeyFile)
	}
	data, err := os.ReadFile(filename)
	if err != nil {
		return errors.Join(err, ErrorParseKeyFile)
	}
	if err := json.Unmarshal(data, kf); err != nil {
		return errors.Join(err, fmt.Errorf("file is not a SASL config (see this example: %s)", ExampleConfig))
	}
	return nil
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

const (
	KeyFileEnvName = "KAFKA_SASL_CREDENTIALS"
	KeyFileName    = "kafka_sasl_credentials.json"
	SecretsDir     = ".secrets"
	ConfigDir      = ".config"
	ConfigSubDir   = "kafka"
	ExampleConfig  = `{
	"key":    "my-api-user",
	"secret": "my-api-secret",
	"server": "my-cluster.europe-west3.gcp.confluent.cloud:9092"
}`
)

func FindKeyFile() (string, error) {
	home := os.Getenv("HOME")
	locations := []string{
		os.Getenv(KeyFileEnvName),                             // 1. use anything on the env
		path.Join(ConfigDir, ConfigSubDir, KeyFileName),       // 2. use a local config
		path.Join(SecretsDir, KeyFileName),                    // 3. use a local secret
		path.Join(home, ConfigDir, ConfigSubDir, KeyFileName), // 4. use the HOME config
		path.Join(home, SecretsDir, KeyFileName),              // 5. use the HOME secret
	}
	var result error
	for _, loc := range locations {
		if f, err := os.Open(loc); result != nil {
			_ = f.Close()
			return f.Name(), nil
		} else {
			result = errors.Join(result, err)
		}
	}
	result = errors.Join(result, errors.New("SASL config not found in the known locations"))
	return "", result
}

func LoadKeyFile() (*KeyFile, error) {
	cfg := KeyFile{}

	if err := cfg.Load(); err != nil {
		return nil, err
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}
