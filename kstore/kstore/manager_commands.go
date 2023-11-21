package kstore

import (
	"context"
	"errors"
)

// Commands
// ========
// The following publish functions are individual commands
// that can be used to manage tables.

// Setup must be called to initialize the SchemaManager
// and setup the TablesInfo topic in Kafka.
func (tm *SchemaManager) Setup(ctx context.Context) error {
	return tm.setup(ctx)
}

// Validate checks the SchemaManager setup for correctness.
func (tm *SchemaManager) Validate() error {
	return tm.validate()
}

// PurgeTable clears the Schema for the given table and deletes the table topic.
func (tm *SchemaManager) PurgeTable(ctx context.Context, schema *tableSchema) error {
	if err := tm.validate(); err != nil {
		return err
	}
	return errors.Join(
		tm.ResetTable(ctx, schema),
		tm.DeleteTable(ctx, schema),
	)
}

// ResetTable clears teh Schema for the given topic.
func (tm *SchemaManager) ResetTable(ctx context.Context, schema *tableSchema) error {
	if err := tm.validate(); err != nil {
		return err
	}
	schema.Schema = nil
	return tm.createOrUpdateTable(ctx, schema)
}

// DeleteTable resets the table schema and deletes the table topic.
func (tm *SchemaManager) DeleteTable(ctx context.Context, schema *tableSchema) error {
	if err := tm.validate(); err != nil {
		return err
	}
	return tm.deleteTable(ctx, schema)
}

// DeleteTopic deletes a topic.
func (tm *SchemaManager) DeleteTopic(ctx context.Context, name string) error {
	if err := tm.validate(); err != nil {
		return err
	}
	return tm.deleteTopics(ctx, name)
}

// CreateOrUpdateTable creates a table topic (if needed) and updates the table schema.
func (tm *SchemaManager) CreateOrUpdateTable(ctx context.Context, schema *tableSchema) error {
	if err := tm.validate(); err != nil {
		return err
	}
	return tm.createOrUpdateTable(ctx, schema)
}

// GetSchema returns the stored schema for the given TableTopic.
func (tm *SchemaManager) GetSchema(ctx context.Context, table TableTopic) (FieldSchema, error) {
	return nil, errors.New("GetSchema not implemented")
}
