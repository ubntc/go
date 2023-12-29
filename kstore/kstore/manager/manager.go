package manager

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/ubntc/go/kstore/kschema"
	"github.com/ubntc/go/kstore/kstore/status"
	"github.com/ubntc/go/kstore/provider/api"
)

type SchemaManager struct {
	schemasTopic string
	client       api.Client
}

func NewSchemaManager(schemasTopic string, client api.Client) *SchemaManager {
	tm := &SchemaManager{
		schemasTopic: schemasTopic,
		client:       client,
	}
	return tm
}

func (tm *SchemaManager) setup(ctx context.Context) error {
	if err := tm.validate(); err != nil {
		return err
	}

	if _, err := tm.createCompactedTopics(ctx, tm.schemasTopic); err != nil {
		return err
	}

	// TODO: manage partitions

	log.Println("initialized SchemaManager with topic:", tm.schemasTopic)
	return nil
}

func (tm *SchemaManager) createOrUpdateTable(ctx context.Context, schema *kschema.Schema) error {
	table := schema.Name
	val, err := json.Marshal(schema)
	if err != nil {
		return err
	}

	msg := kschema.NewMessage(tm.schemasTopic, []byte(table), val)

	topic := schema.GetTopic()

	info, err := tm.createCompactedTopics(ctx, topic)
	if err != nil {
		log.Println("failed to created topic:", topic, ", status:", info.Status(topic))
		return err
	}

	if info.Status(topic) == status.TopicStatusCreated {
		log.Println("created new topic:", topic, " for table:", table)
	}

	err = tm.client.Write(ctx, tm.schemasTopic, msg)
	if err != nil {
		return err
	}
	log.Println("updated table schema:", schema, " for topic:", topic)

	return nil
}

func (tm *SchemaManager) deleteTable(ctx context.Context, schema *kschema.Schema) error {
	msg := kschema.NewMessage(tm.schemasTopic, []byte(schema.Name), nil)

	err := tm.client.Write(ctx, tm.schemasTopic, msg)
	if err != nil {
		return err
	}

	topic := schema.GetTopic()

	err = tm.deleteTopics(ctx, topic)
	if err != nil {
		return err
	}
	log.Println("deteled topic:", topic)

	return nil
}

func (tm *SchemaManager) validate() error {
	switch {
	case tm.client == nil:
		return ErrorWriterNotDefined
	case tm.schemasTopic == "":
		return ErrorEmptyTopic
	default:
		return nil
	}
}

func (tm *SchemaManager) createCompactedTopics(ctx context.Context, topics ...string) (*status.TopicInfo, error) {
	topicErrors, err := tm.client.CreateTopics(ctx, topics...)
	if err != nil {
		return nil, err
	}

	info := &status.TopicInfo{
		Errors:    make(map[string]error),
		StatusMap: make(map[string]status.TopicStatus),
	}

	for name, topicError := range topicErrors {
		switch {
		case topicError == nil:
			info.StatusMap[name] = status.TopicStatusCreated
		case tm.client.IsExistsError(topicError):
			info.StatusMap[name] = status.TopicStatusExists
		default:
			err = errors.Join(err, topicError)
			info.Errors[name] = topicError
		}
	}

	return info, err
}

func (tm *SchemaManager) deleteTopics(ctx context.Context, topics ...string) error {
	_, err := tm.client.DeleteTopics(ctx, topics...)
	if err != nil {
		return err
	}

	return nil
}

func (tm *SchemaManager) Client() api.Client {
	return tm.client
}
