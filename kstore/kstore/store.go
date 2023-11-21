package kstore

import (
	"context"
	"sync"
	"time"
)

type Store struct {
	records map[string][]byte
	table   *tableSchema

	client Client
	mu     sync.RWMutex
}

var StoreAwaitTimeout = time.Second

type Database struct {
	db map[string]*Store

	manager *SchemaManager
	client  Client
	mu      sync.RWMutex
}

func NewDatabase(manager *SchemaManager, client Client) *Database {
	return &Database{
		db:      make(map[string]*Store),
		manager: manager,
		client:  client,
	}
}

func (s *Database) CreateOrUpdateTable(ctx context.Context, table *tableSchema) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	ts, ok := s.db[table.Name]
	if !ok {
		ts = &Store{
			table:   table,
			records: make(map[string][]byte),
			client:  s.client,
		}
	}
	if err := ts.BeginTx(TxWrite, func(ts *Store) error {
		// ensure all messages are compatible with the new schema
		// TODO: only on schema change
		for _, v := range ts.records {
			row := Row{}
			if err := row.Decode(v); err != nil {
				return err
			}
			if err := table.Schema.Validate(row); err != nil {
				return err
			}
		}
		// First send the new/changed schema to the schema topic
		if err := s.manager.createOrUpdateTable(ctx, table); err != nil {
			return err
		}
		// If this is successful then also setup/update the local store
		return nil
	}); err != nil {
		return err
	}
	ts.table.Schema = table.Schema
	s.db[table.Name] = ts
	return nil
}

func (s *Database) awaitGetStore(ctx context.Context, tbl *tableSchema) (*Store, error) {
	ticker := time.NewTicker(StoreAwaitTimeout)
	defer ticker.Stop()
	for {
		store, err := s.GetStore(tbl)
		switch {
		case err == nil:
			return store, nil
		case err != ErrorReadStoreNotInitalized:
			return nil, err
		}
		s.client.GetLogger()("awaiting TableStore:", tbl.Name)
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-ticker.C:
			continue
		}
	}
}

func (s *Database) StartTableReader(ctx context.Context, tbl *tableSchema) (<-chan error, error) {
	r := s.client.NewReader(tbl.GetTopic())
	errch := make(chan error, 1)

	go func() {
		defer close(errch)
		store, err := s.awaitGetStore(ctx, tbl)
		if err != nil {
			errch <- err
			return
		}
		errch <- store.consumeLoop(ctx, r)
	}()

	return errch, nil
}

func (s *Database) GetStore(table *tableSchema) (*Store, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.getStore(table)
}

func (s *Database) getStore(table *tableSchema) (*Store, error) {
	ts, ok := s.db[table.Name]
	if !ok {
		return nil, ErrorStoreNotInitalized
	}
	return ts, nil
}

func (s *Database) WriteRows(ctx context.Context, tbl *tableSchema, rows ...Row) error {
	// do not allow wrting any metadata writes on the stores
	// do allow reading metadata
	// do allow individual store access
	s.mu.RLock()
	defer s.mu.RUnlock()
	ts, err := s.getStore(tbl)
	if err != nil {
		return err
	}
	// starts a write Tx on the individual store
	return ts.BeginTx(TxWrite, func(ts *Store) error {
		return ts.persistRows(ctx, rows...)
	})
}
