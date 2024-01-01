package kafkago

import (
	"context"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/ubntc/go/kstore/kstore"
	"github.com/ubntc/go/kstore/kstore/config"
	"github.com/ubntc/go/kstore/provider/api"
)

type Client struct {
	client        *kafka.Client
	defaultWriter *Writer

	logger     api.LoggerFunc
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
		logger:     kstore.NilLogger(),
		properties: props,
		keyFile:    cfg,
		group:      group,
	}
	// log.Println("creating default writer")
	c.defaultWriter = &Writer{writer: c.newWriter()}
	return c
}

func (c *Client) NewWriter() api.Writer {
	return &Writer{writer: c.newWriter()}
}

func (c *Client) newWriter() *kafka.Writer {
	cfg := c.keyFile
	w := &kafka.Writer{
		Addr:         kafka.TCP(cfg.Server),
		Topic:        "", // leave topic empty, must be set when writing messages
		Transport:    defaultTransport(cfg.Key, cfg.Secret),
		BatchSize:    1,
		Logger:       kafka.LoggerFunc(c.logger),
		Balancer:     kafka.Murmur2Balancer{},
		WriteTimeout: time.Second,
		Async:        false,
		// MaxAttempts:  1,
	}
	log.Printf("creating writer for server: %v", w.Addr)
	return w
}

func (c *Client) NewReader(topic string) api.Reader {
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
	log.Printf("creating reader for topic: %s (group:%v)", r.topic, c.group)
	return r
}

func (c *Client) CreateTopics(ctx context.Context, topics ...string) (api.TopicErrors, error) {
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

func (c *Client) DeleteTopics(ctx context.Context, topics ...string) (api.TopicErrors, error) {
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
func (c *Client) Write(ctx context.Context, topic string, msg ...api.Message) error {
	return c.defaultWriter.Write(ctx, topic, msg...)
}

func (c *Client) Read(ctx context.Context, topic string, partition int, offset *uint64) (api.Message, error) {
	r := c.NewReader(topic).(*Reader)
	defer r.Close()

	cfg := r.reader.Config()
	cfg.GroupID = ""
	cfg.Partition = partition
	cfg.StartOffset = kafka.FirstOffset
	r.reader = kafka.NewReader(cfg)

	if offset != nil {
		if err := r.reader.SetOffset(int64(*offset)); err != nil {
			return nil, err
		}
		if r.reader.Offset() < 0 {
			return nil, ErrInvalidOffset
		}
	}

	log.Printf("changed reader config to group='' and partition=%d to read from offset=%d", partition, offset)
	return r.Read(ctx)
}

func (c *Client) SetLogger(fn api.LoggerFunc) {
	c.logger = fn
}

func (c *Client) GetLogger() api.LoggerFunc {
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
var _ = api.Client(&Client{})
