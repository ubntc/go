package pebble

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"sync"
	"time"

	"github.com/cockroachdb/pebble"
	"github.com/ubntc/go/kstore/kschema"
	"github.com/ubntc/go/kstore/provider/api"
)

type AcquireMode int

const (
	AcquireModeRead AcquireMode = iota
	AcquireModeWrite
)

type Client struct {
	db     map[string]*pebble.DB
	logger api.LoggerFunc
	prefix string

	mu sync.RWMutex
}

type HandleFunc func(ctx context.Context, message *Message) error

func NewClient(prefix string) *Client {
	switch prefix {
	case "":
		log.Println("using os.TempDir as pebble dir")
		prefix = path.Join(os.TempDir(), "pebble")
	case ".pebble":
		log.Println("using .pebble as pebble dir")
	default:
		log.Printf("using %s/.pebble as pebble dir", prefix)
		prefix = path.Join(prefix, ".pebble")
	}
	c := &Client{
		db:     make(map[string]*pebble.DB),
		prefix: prefix,
	}

	return c
}

func (c *Client) NewWriter() api.Writer {
	return NewWriter(c)
}

func (c *Client) NewReader(topic string) api.Reader {
	log.Printf("creating reader for pebble topic: %s\n", topic)
	r := NewReader(c, topic, StartOffsetFirst)
	return r
}

func (c *Client) CreateTopics(ctx context.Context, topics ...string) (api.TopicErrors, error) {
	result := make(map[string]error)
	log.Printf("creating %d topics as pebble tables", len(topics))
	for _, t := range topics {
		if createErr := c.CreateDB(t); createErr != nil {
			log.Println(createErr)
			result[t] = createErr
		}
	}
	if len(result) > 0 {
		return result, fmt.Errorf("failed to create %d of %d pebble topics", len(result), len(topics))
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
		return result, fmt.Errorf("failed to delete %d of %d pebble topics", len(result), len(topics))
	}
	return nil, nil
}

// Write writes a message.
func (c *Client) Write(ctx context.Context, topic string, msg ...api.Message) error {
	db, release, err := c.AcquireDB(topic, AcquireModeWrite)
	if err != nil {
		return err
	}
	defer release()

	for _, m := range msg {
		sk := StorageKey(m)
		// log.Printf("writing message: %s with storageKey: %v", m.String(), sk)
		if err := db.Set(sk, StorageValue(m), pebble.Sync); err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) Read(ctx context.Context, topic string, partition int, offset *uint64) (api.Message, error) {
	if offset == nil {
		return c.ReadNext(ctx, topic, nil)
	}
	return c.ReadNext(ctx, topic, OffsetBytes(*offset))
}

func (c *Client) Get(topic string, storageKey []byte) (api.Message, error) {
	db, release, err := c.AcquireDB(topic, AcquireModeRead)
	if err != nil {
		return nil, err
	}
	defer release()
	value, closer, err := db.Get(storageKey)
	if err != nil {
		return nil, err
	}
	defer closer.Close()

	m := &Message{kschema.Message{}}
	if err := m.Decode(value); err != nil {
		return nil, err
	}

	// log.Println("read message:", m.String())

	return m, nil
}

func (c *Client) ReadNext(ctx context.Context, topic string, currentKey []byte) (api.Message, error) {
	ticker := time.NewTicker(time.Millisecond * 100)
	defer ticker.Stop()

	var key, value []byte
	next := func() (foundNext bool) {
		db, release, err := c.AcquireDB(topic, AcquireModeRead)
		if err != nil {
			return false
		}
		defer release()

		// logMsg := "Read.next:"

		opts := pebble.IterOptions{}
		if currentKey != nil {
			// logMsg += fmt.Sprintf(" 🔑=%v", Offset(currentKey))
			opts.LowerBound = currentKey
		}

		iter := db.NewIterWithContext(ctx, &opts)
		defer iter.Close()

		// set the iterator to the first position
		if !iter.First() {
			// there is no new message yet
			return false
		}
		// We have a valid iterator now, which means:
		// A) opts.LowerBound is NOT set -> iterator is at the desired next position
		// B) opts.LowerBound IS set     -> iterator is at the last known position
		//    The desired position should be the next one.

		// Case A: LowerBound not set -> return the first message
		if opts.LowerBound == nil {
			key = iter.Key()
			value = iter.Value()
			// log.Println(logMsg, "📬 read.first:", "key", key, "value", string(value))
			return true
		}

		// Case B: LowerBound is set -> return the next message
		if iter.Next() {
			key = iter.Key()
			value = iter.Value()
			// log.Println(logMsg, "📬 read.next:", "key", key, "value", string(value))
			return true
		}

		// there is no new message yet
		return false
	}

	// read forever until we get a message or the context is done
	for !next() {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-ticker.C:
			continue
		}
	}

	msg := kschema.RawMessage(topic, Offset(key), key, value)

	return &msg, nil
}

func (c *Client) FindFirst(ctx context.Context, topic string) ([]byte, error) {
	db, release, err := c.AcquireDB(topic, AcquireModeRead)
	if err != nil {
		return nil, errors.Join(err, fmt.Errorf("failed to acquire table for find first: %s", topic))
	}
	defer release()

	iter := db.NewIter(nil)
	defer iter.Close()
	if iter.First() {
		return iter.Key(), nil
	}
	return nil, nil
}

func (c *Client) FindLast(ctx context.Context, topic string) ([]byte, error) {
	db, release, err := c.AcquireDB(topic, AcquireModeRead)
	if err != nil {
		return nil, errors.Join(err, fmt.Errorf("failed to acquire table for find last: %s", topic))
	}
	defer release()

	iter := db.NewIter(nil)
	defer iter.Close()
	if iter.Last() {
		return iter.Key(), nil
	}
	return nil, nil
}

func (c *Client) Subscribe(ctx context.Context, topic string, fn HandleFunc) error {
	return nil
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

func (c *Client) AcquireDB(topic string, mode AcquireMode) (*pebble.DB, context.CancelFunc, error) {
	var release context.CancelFunc

	if mode == AcquireModeRead {
		c.mu.RLock()
		release = c.mu.RUnlock
	} else {
		c.mu.Lock()
		release = c.mu.Unlock
	}

	db, ok := c.db[topic]
	if !ok {
		release()
		return nil, nil, ErrorTableNotInitalized
	}
	return db, release, nil
}

func (c *Client) CreateDB(topic string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.db[topic]; ok {
		return nil
	}

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
		log.Println("close error:", err)
		return err
	}

	if err := os.RemoveAll(c.dbPath(topic)); err != nil {
		return err
	}

	return nil
}

func (c *Client) dbPath(topic string) string {
	return path.Join(c.prefix, topic)
}

var _ = api.Client(&Client{})
