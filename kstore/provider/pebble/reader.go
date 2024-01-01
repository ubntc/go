package pebble

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/ubntc/go/kstore/provider/api"
)

type StartOffset int

const (
	StartOffsetFirst StartOffset = iota
	StartOffsetLast
)

type Reader struct {
	topic       string
	client      *Client
	startOffset StartOffset

	lastCommittedStorageKey []byte // last committed key
	// TODO: persist reader keys per group ID in pebble

	mu sync.RWMutex
}

func NewReader(client *Client, topic string, startOffset StartOffset) *Reader {
	r := &Reader{
		topic:       topic,
		client:      client,
		startOffset: startOffset,
	}
	if err := r.Validate(); err != nil {
		panic(err)
	}
	return r
}

func (r *Reader) Validate() error {
	var result error
	if r.topic == "" {
		result = errors.Join(result, fmt.Errorf("topic not defined"))
	}
	if r.client == nil {
		result = errors.Join(result, fmt.Errorf("client not defined"))
	}
	return result
}

// recoverLastCommit initializes the reader by setting the currentStorageKey
// to the last committed key found in the DB.
//
// NOTE: Must be protected by r.mu!
func (r *Reader) recoverLastCommit(ctx context.Context) error {
	if r.lastCommittedStorageKey != nil {
		return nil
	}

	var key []byte
	var err error
	switch r.startOffset {
	case StartOffsetLast:
		key, err = r.client.FindLast(ctx, r.topic)
	case StartOffsetFirst:
		key, err = r.client.FindFirst(ctx, r.topic)
	default:
		return ErrorInvalidStartOffset
	}

	if err != nil {
		return errors.Join(err, ErrorOffsetNotFound)
	}

	r.lastCommittedStorageKey = key
	log.Println("initalized", r.describe())
	return nil
}

func (r *Reader) key() []byte {
	return r.lastCommittedStorageKey
}

// describe describes the current state of the reader.
//
// NOTE: Must be protected by r.mu!
func (r *Reader) describe() string {
	reader := "empty reader"
	status := ""
	if r.lastCommittedStorageKey != nil {
		reader = "reader"
		status = fmt.Sprintf("at offset=%d", Offset(r.key()))
	}
	return fmt.Sprintf("%s for topic=%s %s", reader, r.topic, status)
}

func (r *Reader) Read(ctx context.Context) (api.Message, error) {
	// Make sure the reader state is not change while we read.
	// We need a write lock here, since we need to initialize the reader on the first read.
	r.mu.Lock()
	defer r.mu.Unlock()

	// in-memory buffer not defined -> read from the disk
	return r.readFromDB(ctx)
}

func (r *Reader) readFromDB(ctx context.Context) (api.Message, error) {
	// lazy init reader
	if err := r.recoverLastCommit(ctx); err != nil {
		return nil, err
	}

	msg, err := r.client.ReadNext(ctx, r.topic, r.lastCommittedStorageKey)
	if err != nil {
		return nil, err
	}

	status := CompareOffsetByKey(r.key(), StorageKey(msg))
	Metrics.ObserveRead(msg, r.topic, status)

	if status < OffsetStatusCurrent {
		return nil, ErrorReicevedOldMessage
	}

	return msg, nil
}

func (r *Reader) Commit(ctx context.Context, msg api.Message) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if CompareOffsetByKey(r.lastCommittedStorageKey, StorageKey(msg)) < OffsetStatusCurrent {
		panic("message already seen, cannot commit old messages")
	}

	r.lastCommittedStorageKey = StorageKey(msg)
	return nil
}

func (r *Reader) Get(storageKey []byte) (api.Message, error) {
	return r.client.Get(r.topic, storageKey)
}

func (r *Reader) Close() error {
	return nil
}

// ensure we implement the full interface
func init() { _ = api.Reader(&Reader{}) }
