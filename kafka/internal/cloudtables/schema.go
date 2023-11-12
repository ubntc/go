package cloudtables

type FieldType string

const (
	FieldTypeString  = "string"
	FieldTypeInt64   = "int64"
	FieldTypeFloat64 = "float64"
	FieldTypeBool    = "bool"
	FieldTypeRecord  = "record"
)

type Field struct {
	Name       string    `json:"name,omitempty"`
	Type       FieldType `json:"type,omitempty"`
	Repeated   bool      `json:"repeated,omitempty"`
	RecordType Schema    `json:"record_type,omitempty"`
}

type Schema []Field

type Table struct {
	Name   string `json:"name,omitempty"`
	Schema Schema `json:"schema,omitempty"`
}

type TableState struct {
	Table  `json:"table,omitempty"`
	Status string `json:"status,omitempty"`
}
