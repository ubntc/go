package kstore

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"

	"github.com/ubntc/go/kstore/kschema"
	"github.com/ubntc/go/kstore/provider/api"
)

type Store struct {
	records map[string][]byte
	table   *kschema.Schema

	client api.Client
	mu     sync.RWMutex
}

var StoreAwaitTimeout = time.Second

func (s *Store) consumeLoop(ctx context.Context, reader api.Reader) (consumeErr error) {
	defer reader.Close()
	defer func() { consumeErr = FilterGraceful(consumeErr) }()

	storeAndCommit := func(m api.Message) (result error) {
		s.mu.Lock()
		// commit after unlocking
		defer func() {
			if result == nil { // commit only on successful storage
				result = errors.Join(result, reader.Commit(ctx, m))
			}
		}()
		// unlock directly after applying the change and before committing
		// it is safe to see an uncommitted message again, once we have a version check
		defer s.mu.Unlock()
		// store the new message
		// TODO: check version and reject store if it is too old
		return s.storeMessages(ctx, m)
	}

	log.Println("starting consumeLoop for topic:", s.table.GetTopic())
	defer log.Println("stopped consumeLoop for topic:", s.table.GetTopic())
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		m, err := reader.Read(ctx)
		if err != nil {
			return err
		}

		row := kschema.Row{Key: m.Key()}
		if err := row.Decode(m.Value()); err != nil {
			return
		}

		if err := storeAndCommit(m); err != nil {
			return
		}
	}
}

func (s *Store) storeMessages(ctx context.Context, values ...api.Message) error {
	log.Printf("storeMessages: %d rows\n", len(values))
	for _, m := range values {
		s.records[string(m.Key())] = m.Value()
	}
	return nil
}

func (s *Store) persistRows(ctx context.Context, rows ...kschema.Row) (result error) {
	// ensure to write data locally after successfully sending it out
	log.Printf("persistRows: %d rows\n", len(rows))
	topic := s.table.GetTopic()

	var err error
	var rowBytes []byte
	var messages []api.Message

	// store committed messages and augment the returned error
	defer func() {
		result = errors.Join(
			result,
			s.storeMessages(ctx, messages...),
		)
	}()

	for _, r := range rows {
		if err = s.table.Schema.Validate(r); err != nil {
			return err
		}
		if rowBytes, err = r.Encode(); err != nil {
			return err
		}
		msg := kschema.NewMessage(topic, r.Key, rowBytes)
		err = s.client.Write(ctx, topic, msg)
		if err != nil {
			return err
		}
		messages = append(messages, msg)
	}

	return nil
}

// ---------------------------
// Transaction Processing,
// Locked Reads, Locked Writes
// ---------------------------

type (
	TxFunc func(ks *Store) error
	TxType string
)

const (
	TxRead  = "read"
	TxWrite = "write"
)

func (ts *Store) BeginTx(typ TxType, fn TxFunc) error {
	if typ == TxWrite {
		ts.mu.Lock()
		defer ts.mu.Unlock()
	} else {
		ts.mu.RLock()
		defer ts.mu.RUnlock()
	}
	log.Println("BeginTx:", typ)
	return fn(ts)
}

func (ts *Store) WriteRow(ctx context.Context, value kschema.Row) error {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	return ts.persistRows(ctx, value)
}

func (ts *Store) GetRow(ctx context.Context, key string) (*kschema.Row, error) {
	ts.mu.RLock()
	defer ts.mu.RUnlock()
	data, ok := ts.records[key]
	if !ok {
		return nil, nil
	}
	return (&kschema.Row{}).Decoded(data)
}
