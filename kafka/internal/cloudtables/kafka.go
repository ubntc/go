package cloudtables

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/plain"
)

var (
	ErrorNotInitialized   = errors.New("TableManager not initialized")
	ErrorWriterNotDefined = errors.New("Writer not defined")
	ErrorTopicMismatch    = errors.New("Writer topic does not match manager topic")
)

func NewLogger(name string) kafka.LoggerFunc {
	return kafka.LoggerFunc(func(format string, args ...interface{}) {
		log.Printf(name+": "+format, args...)
	})
}

func NilLogger() kafka.LoggerFunc {
	return kafka.LoggerFunc(func(format string, args ...interface{}) {})
}

type TableManager struct {
	topic *kafka.Topic

	Topic  string
	Writer *kafka.Writer
	Reader *kafka.Reader
}

// Setup must be called to initialize the TableManager
// and setup the TablesInfo topic in Kafka.
func (tm *TableManager) Setup(ctx context.Context) error {
	tm.topic = &kafka.Topic{
		Name: tm.Topic,
	}

	if tm.Writer == nil {
		return ErrorWriterNotDefined
	}

	if tm.Writer.Topic == "" {
		tm.Writer.Topic = tm.Topic
		log.Println("setting unset Writer topic to Manager topic:", tm.Topic)
	}

	if tm.Writer.Topic != tm.Topic {
		return ErrorTopicMismatch
	}

	if _, err := tm.createCompactedTopic(ctx, tm.Topic); err != nil {
		return err
	}

	// TODO: setup partitions

	log.Println("TableManager initialized with topic:", tm.Topic)

	return nil
}

func (tm *TableManager) CreateOrUpdateTable(ctx context.Context, cfg Table) error {
	if tm.topic == nil {
		return ErrorNotInitialized
	}

	table := cfg.Name
	val, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	msg := kafka.Message{
		Key:   []byte(table),
		Value: val,
	}

	code, err := tm.createCompactedTopic(ctx, tm.TopicForTable(cfg.Name))
	if err != nil && code != kafka.TopicAlreadyExists {
		return err
	}

	err = tm.Writer.WriteMessages(ctx, msg)
	if err != nil {
		return err
	}

	if code == kafka.TopicAlreadyExists {
		log.Println("updated table schema:", cfg)
	}
	log.Println("created table schema:", cfg)

	return nil
}

func (tm *TableManager) DeleteTable(ctx context.Context, name string) error {
	if tm.topic == nil {
		return ErrorNotInitialized
	}

	msg := kafka.Message{
		Key:   []byte(name),
		Value: nil,
	}
	err := tm.Writer.WriteMessages(ctx, msg)
	if err != nil {
		return err
	}

	topic := tm.TopicForTable(name)

	err = tm.deleteTopic(ctx, topic)
	if err != nil {
		return err
	}
	log.Println("deteled topic:", topic)

	return nil
}

func (tm *TableManager) createCompactedTopic(ctx context.Context, topic string) (kafka.Error, error) {
	c := kafka.Client{
		Addr:      tm.Writer.Addr,
		Transport: tm.Writer.Transport,
	}

	configEntries := []kafka.ConfigEntry{}

	for k, v := range DefaultCompactConfig() {
		configEntries = append(configEntries, kafka.ConfigEntry{ConfigName: k, ConfigValue: v})
	}

	req := &kafka.CreateTopicsRequest{
		Addr: c.Addr,
		Topics: []kafka.TopicConfig{
			{
				Topic:             topic,
				NumPartitions:     DefaultNumPartitions,
				ReplicationFactor: DefaultReplicationFactor,
				ConfigEntries:     configEntries,
			},
		},
	}

	res, err := c.CreateTopics(ctx, req)
	if err != nil {
		return kafka.Unknown, err
	}

	for name, topicError := range res.Errors {
		switch ErrorCode(topicError) {
		case 0:
			log.Println("topic created:", name)
		case kafka.TopicAlreadyExists:
			log.Println("topic exists:", name)
		default:
			log.Println("failed to create topic:", name)
			err = errors.Join(err, topicError)
		}
	}

	if err != nil {
		return kafka.Unknown, err
	}

	return 0, nil
}

func (tm *TableManager) deleteTopic(ctx context.Context, topic string) error {
	c := kafka.Client{
		Addr:      tm.Writer.Addr,
		Transport: tm.Writer.Transport,
	}

	_, err := c.DeleteTopics(ctx, &kafka.DeleteTopicsRequest{
		Addr:   tm.Writer.Addr,
		Topics: []string{topic},
	})
	if err != nil {
		return err
	}

	return nil
}

func (tn *TableManager) TopicForTable(name string) string {
	return DefaultPrefix + name
}

func (tn *TableManager) TableForTopic(name string) string {
	if v, ok := strings.CutPrefix(name, DefaultPrefix); ok {
		return v
	}
	return ""
}

var reErrorCode = regexp.MustCompile(`^\[[0-9]+\]`)

func ErrorCode(err error) kafka.Error {
	if err == nil {
		return 0
	}

	match := reErrorCode.FindString(err.Error())
	if len(match) < 3 {
		return kafka.Unknown
	}
	code, _ := strconv.Atoi(match[1 : len(match)-1])
	return kafka.Error(code)
}

func DefaultTransport(cfg *KeyFile) *kafka.Transport {
	transport := &kafka.Transport{
		SASL: plain.Mechanism{
			Username: cfg.Key,
			Password: cfg.Secret,
		},
		TLS: &tls.Config{
			MinVersion: tls.VersionTLS12,
		},
		DialTimeout: time.Second * 20,
	}
	return transport
}

func NewWriter(cfg *KeyFile) *kafka.Writer {
	return &kafka.Writer{
		Addr:      kafka.TCP(cfg.Server),
		Topic:     DefaultManagerTopic,
		Transport: DefaultTransport(cfg),
		BatchSize: 1,
		Logger:    NilLogger(),
		Balancer:  kafka.Murmur2Balancer{},
	}
}
