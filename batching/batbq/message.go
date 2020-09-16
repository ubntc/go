package batbq

import (
	"log"

	"cloud.google.com/go/bigquery"
)

// Message defines an (n)ackable message that contains the data for BigQuery.
type Message interface {
	Data() *bigquery.StructSaver
	Ack()
	Nack(err error)
}

// LogMessage implements the `Message` interface. A LogMessage
// ignores the `Ack()` and logs a given error from `Nack(err error)`.
type LogMessage struct {
	bigquery.StructSaver
}

// Ack does nothing.
func (m *LogMessage) Ack() {}

// Nack logs the error.
func (m *LogMessage) Nack(err error) {
	log.Printf("LogMessage Nacked with error: %v, data: %v", err, m.Data())
}

// Data returns the embedded StructSaver.
func (m *LogMessage) Data() *bigquery.StructSaver {
	return &m.StructSaver
}

func toStructs(messages []Message) []*bigquery.StructSaver {
	res := make([]*bigquery.StructSaver, len(messages))
	for i, m := range messages {
		res[i] = m.Data()
	}
	return res
}
