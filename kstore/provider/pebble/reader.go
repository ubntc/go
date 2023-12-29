package pebble

import (
	"bytes"
	"context"
	"sync"
	"time"

	"github.com/cockroachdb/pebble"
	"github.com/ubntc/go/kstore/kschema"
	"github.com/ubntc/go/kstore/provider/api"
)

type Reader struct {
	topic  string
	client *Client

	key []byte // last committed key
	// TODO: persist reader keys per group ID in pebble

	mu sync.Mutex
}

type StartOffset int

const (
	StartOffsetFirst StartOffset = iota
	StartOffsetLast
)

func NewReader(client *Client, topic string, offset StartOffset) (*Reader, error) {
	r := &Reader{
		topic:  topic,
		client: client,
	}

	db, err := r.client.GetDB(topic)
	if err != nil {
		return nil, err
	}

	iter := db.NewIter(nil)
	switch offset {
	case StartOffsetFirst:
		iter.First()
	case StartOffsetLast:
		iter.Last()
	default:
		iter.First()
	}

	r.key = iter.Key()
	return r, nil
}

func (r *Reader) Close() error {
	return nil
}

func (r *Reader) Commit(ctx context.Context, msg api.Message) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.key = StorageKey(msg)
	return nil
}

func (r *Reader) Read(ctx context.Context) (api.Message, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	db, err := r.client.GetDB(r.topic)
	if err != nil {
		return nil, err
	}

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	var key, value []byte
	next := func() bool {
		iter := db.NewIterWithContext(ctx, &pebble.IterOptions{LowerBound: r.key})
		defer iter.Close()
		for iter.Next() {
			if bytes.Equal(iter.Key(), r.key) {
				// skip until we find a valid new key
				continue
			}
			key = iter.Key()
			value = iter.Value()
			return true
		}
		return false
	}

	for !next() {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-ticker.C:
			continue
		}
	}

	msg := kschema.RawMessage(r.topic, Offset(key), key, value)
	return &msg, nil
}

// ensure we implement the full interface
func init() { _ = api.Reader(&Reader{}) }
