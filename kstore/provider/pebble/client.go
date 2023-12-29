package pebble

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"sync"

	"github.com/cockroachdb/pebble"
	"github.com/ubntc/go/kstore/provider/api"
)

const prefix = ".pebble"

type Client struct {
	db            map[string]*pebble.DB
	defaultWriter *Writer
	logger        api.LoggerFunc

	mu sync.RWMutex
}

func NewClient() *Client {
	c := &Client{
		db: make(map[string]*pebble.DB),
	}

	log.Println("creating default writer")
	c.defaultWriter = NewWriter(c)
	return c
}

func (c *Client) NewWriter() api.Writer {
	return NewWriter(c)
}

func (c *Client) NewReader(topic string) api.Reader {
	log.Printf("creating reader for topic: %s\n", topic)
	return nil
}

func (c *Client) CreateTopics(ctx context.Context, topics ...string) (api.TopicErrors, error) {
	result := make(map[string]error)
	log.Printf("creating %d topics", len(topics))
	for _, t := range topics {
		if createErr := c.CreateDB(t); createErr != nil {
			result[t] = createErr
		}
	}
	if len(result) > 0 {
		return result, fmt.Errorf("failed to create %d of %d topics", len(result), len(topics))
	}
	return nil, nil
}

func (c *Client) DeleteTopics(ctx context.Context, topics ...string) (api.TopicErrors, error) {
	result := make(map[string]error)
	log.Printf("deleting topics: %v", topics)
	for _, t := range topics {
		if err := c.DeleteDB(t); err != nil {
			result[t] = err
		}
	}
	if len(result) > 0 {
		return result, fmt.Errorf("failed to delete %d of %d topics", len(result), len(topics))
	}
	return nil, nil
}

// Write writes a message using the default writer.
func (c *Client) Write(ctx context.Context, topic string, msg ...api.Message) error {
	return c.defaultWriter.Write(ctx, topic, msg...)
}

func (c *Client) SetLogger(fn api.LoggerFunc) {
	c.logger = fn
}

func (c *Client) GetLogger() api.LoggerFunc {
	return c.logger
}

func (c *Client) IsExistsError(err error) bool {
	return errors.Is(err, pebble.ErrDBAlreadyExists)
}

func (c *Client) Close() error {
	return nil
}

func (c *Client) GetDB(topic string) (*pebble.DB, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	db, ok := c.db[topic]
	if !ok {
		return nil, fmt.Errorf("pebble.DB for topic %s not initalized", topic)
	}
	return db, nil
}

func (c *Client) CreateDB(topic string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	db, err := pebble.Open(c.dbPath(topic), &pebble.Options{})
	if err != nil {
		return err
	}
	c.db[topic] = db
	return nil
}

func (c *Client) closeDB(topic string) error {
	if db, ok := c.db[topic]; ok {
		if err := db.Close(); err != nil {
			return err
		}
		delete(c.db, topic)
	}
	return nil
}

func (c *Client) DeleteDB(topic string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if err := c.closeDB(topic); err != nil {
		return err
	}

	if err := os.RemoveAll(c.dbPath(topic)); err != nil {
		return err
	}

	return nil
}

func (c *Client) dbPath(topic string) string {
	return path.Join(prefix, topic)
}

func dummy() {
	db, err := pebble.Open("demo", &pebble.Options{})
	if err != nil {
		log.Fatal(err)
	}
	key := []byte("hello")
	if err := db.Set(key, []byte("world"), pebble.Sync); err != nil {
		log.Fatal(err)
	}
	value, closer, err := db.Get(key)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s %s\n", key, value)
	if err := closer.Close(); err != nil {
		log.Fatal(err)
	}
	if err := db.Close(); err != nil {
		log.Fatal(err)
	}
}

var _ = api.Client(&Client{})
