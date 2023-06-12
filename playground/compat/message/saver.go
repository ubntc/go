package message

import (
	"cloud.google.com/go/bigquery"
	"github.com/google/uuid"
)

func (m *Msg) Save() (row map[string]bigquery.Value, insertID string, err error) {
	return map[string]bigquery.Value{
			"type":  bigquery.Value(m.Type),
			"value": bigquery.Value(m.Value),
		},
		uuid.NewString(),
		nil
}
