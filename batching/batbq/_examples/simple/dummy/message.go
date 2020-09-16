package dummy

import "cloud.google.com/go/bigquery"

// Message is a dummy message.
type Message struct {
	ID  string
	Val int
}

// ConfirmMessage does nothing.
func (m *Message) ConfirmMessage() {}

// Save implements the ValueSaver interface.
func (m *Message) Save() (row map[string]bigquery.Value, insertID string, err error) {
	v := bigquery.Value(m.Val)
	return map[string]bigquery.Value{"val": v}, m.ID, nil
}
