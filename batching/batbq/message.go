package batbq

import (
	"log"

	"cloud.google.com/go/bigquery"
)

// Message defines an (n)ackable message that contains the data for BigQuery.
type Message interface {
	Data() bigquery.ValueSaver // Data returns a ValueSaver for the bigquery.Inserter
	Ack()                      // Ack confirms successful processing of the message at the sender.
	Nack(err error)            // Nack reports unsuccessful processing and errors to the sender.
}

// LogMessage implements the `Message` interface. A LogMessage ignores the `Ack()` and logs a given
// error from `Nack(err error)`. Use it for testing and for naive data pipelines that do not require
// acknowledging messages.
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
func (m *LogMessage) Data() bigquery.ValueSaver {
	return &m.StructSaver
}
