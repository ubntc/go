package kschema

import (
	"encoding/json"
)

type Row struct {
	Key    []byte `json:"key,omitempty"`
	Values []any  `json:"values,omitempty"`
}

func (r *Row) Encode() ([]byte, error) {
	return json.Marshal(r)
}

func (r *Row) Decode(data []byte) error {
	return json.Unmarshal(data, r)
}

func (r *Row) Decoded(data []byte) (*Row, error) {
	err := json.Unmarshal(data, r)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (r *Row) GetKey() []byte {
	if r == nil {
		return nil
	}
	return r.Key
}

func (r *Row) GetValues() []any {
	if r == nil {
		return nil
	}
	return r.Values
}

type Codec interface {
	GetKey() []byte           // GetKey gets the key bytes
	Encode() ([]byte, error)  // Encode encodes the entire record as bytes (incl. the key)
	Decode(data []byte) error // Decode decodes the data bytes into the record (incl. the key)
	// Row() ([]any, error)      // Row returns the record as kstore.Row
}
