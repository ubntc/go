package kafkago

import (
	"crypto/tls"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/plain"
	"github.com/ubntc/go/kstore/kstore"
	"github.com/ubntc/go/kstore/kstore/config"
)

func NewDialer(cfg *config.KeyFile) *kafka.Dialer {
	transport := defaultTransport(cfg.Key, cfg.Secret)
	return &kafka.Dialer{
		DialFunc:        transport.Dial,
		SASLMechanism:   transport.SASL,
		TLS:             transport.TLS,
		Timeout:         transport.DialTimeout,
		TransactionalID: "", // TODO: evaluate if transactions can be used for "table transactions".
	}
}

func defaultTransport(key, secret string) *kafka.Transport {
	transport := &kafka.Transport{
		SASL: plain.Mechanism{
			Username: key,
			Password: secret,
		},
		TLS: &tls.Config{
			MinVersion: tls.VersionTLS12,
		},
		DialTimeout: time.Second * 20,
	}
	return transport
}

func readerConfig(cfg *config.KeyFile, topic string, group config.Group) kafka.ReaderConfig {
	dialer := NewDialer(cfg)
	return kafka.ReaderConfig{
		Brokers:        []string{cfg.Server},
		GroupID:        group.ID,
		Topic:          topic,
		GroupTopics:    group.Topics, // used only for group management
		Dialer:         dialer,
		Logger:         kafka.LoggerFunc(kstore.NilLogger()),
		IsolationLevel: kafka.ReadCommitted,
		StartOffset:    kafka.FirstOffset, // start from the beginning for new groupID
		CommitInterval: 0,                 // use sync commits
	}
}

func ConfigEntries(props config.KafkaProperties) []kafka.ConfigEntry {
	configEntries := make([]kafka.ConfigEntry, len(props))
	for k, v := range props {
		configEntries = append(configEntries, kafka.ConfigEntry{ConfigName: k, ConfigValue: v})
	}
	return configEntries
}

func DefaultTopicConfigs(props config.KafkaProperties, topics ...string) []kafka.TopicConfig {
	topicConfigs := []kafka.TopicConfig{}
	for _, t := range topics {
		topicConfigs = append(topicConfigs, kafka.TopicConfig{
			Topic:             t,
			NumPartitions:     config.DefaultNumPartitions,
			ReplicationFactor: config.DefaultReplicationFactor,
			ConfigEntries:     ConfigEntries(props),
		})
	}
	return topicConfigs
}
