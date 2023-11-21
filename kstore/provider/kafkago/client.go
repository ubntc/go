package kafkago

import (
	"context"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/ubntc/go/kstore/kstore"
	"github.com/ubntc/go/kstore/kstore/config"
)

type Client struct {
	client        *kafka.Client
	defaultWriter *Writer

	logger     kstore.LoggerFunc
	properties map[string]string
	keyFile    *config.KeyFile
	group      config.Group
}

func NewClient(cfg *config.KeyFile, props config.KafkaProperties, group config.Group) *Client {
	kc := &kafka.Client{
		Addr:      kafka.TCP(cfg.Server),
		Transport: defaultTransport(cfg.Key, cfg.Secret),
	}
	c := &Client{
		client:     kc,
		logger:     NilLogger(),
		properties: props,
		keyFile:    cfg,
		group:      group,
	}
	log.Println("creating default writer")
	c.defaultWriter = &Writer{writer: c.newWriter()}
	return c
}

func (c *Client) NewWriter() kstore.Writer {
	return &Writer{writer: c.newWriter()}
}

func (c *Client) newWriter() *kafka.Writer {
	cfg := c.keyFile
	w := &kafka.Writer{
		Addr:         kafka.TCP(cfg.Server),
		Topic:        "", // leave topic empty, must be set as Message.Topic
		Transport:    defaultTransport(cfg.Key, cfg.Secret),
		BatchSize:    1,
		Logger:       kafka.LoggerFunc(c.logger),
		Balancer:     kafka.Murmur2Balancer{},
		WriteTimeout: time.Second,
		Async:        false,
		// MaxAttempts:  1,
	}
	log.Printf("creating writer for server: %v\n", w.Addr)
	return w
}

func (c *Client) NewReader(topic string) kstore.Reader {
	topics := []string{topic}
	for _, v := range c.group.Topics {
		if v != topic {
			topics = append(topics, v)
		}
	}
	cfg := readerConfig(c.keyFile, topic, config.Group{
		ID:     c.group.ID,
		Topics: topics,
	})
	cfg.Logger = kafka.LoggerFunc(c.logger)

	r := &Reader{
		topic:  topic,
		reader: kafka.NewReader(cfg),
	}
	log.Printf("creating reader for topic: %s (group:%v)\n", r.topic, c.group)
	return r
}

func (c *Client) CreateTopics(ctx context.Context, topics ...string) (kstore.TopicErrors, error) {
	req := &kafka.CreateTopicsRequest{
		Addr:   c.client.Addr,
		Topics: DefaultTopicConfigs(c.properties, topics...),
	}
	log.Printf("creating %d topics", len(req.Topics))
	res, err := c.client.CreateTopics(ctx, req)
	if err != nil {
		return res.Errors, err
	}
	return nil, nil
}

func (c *Client) DeleteTopics(ctx context.Context, topics ...string) (kstore.TopicErrors, error) {
	req := &kafka.DeleteTopicsRequest{
		Addr:   c.client.Addr,
		Topics: topics,
	}
	log.Printf("deleting topics: %v", req.Topics)
	res, err := c.client.DeleteTopics(ctx, req)
	if err != nil {
		return res.Errors, err
	}
	return nil, nil
}

// Write writes a message using the default writer.
func (c *Client) Write(ctx context.Context, topic string, msg ...kstore.Message) error {
	return c.defaultWriter.Write(ctx, topic, msg...)
}

func (c *Client) SetLogger(fn kstore.LoggerFunc) {
	c.logger = fn
}

func (c *Client) GetLogger() kstore.LoggerFunc {
	return c.logger
}

func (c *Client) IsExistsError(err error) bool {
	if err := KafkaError(err); err == kafka.TopicAlreadyExists {
		return true
	}
	return false
}

func (c *Client) Close() error {
	return c.defaultWriter.Close()
}

// ensure we implement the full interface
func init() { _ = kstore.Client(&Client{}) }
