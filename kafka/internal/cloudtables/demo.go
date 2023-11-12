package cloudtables

import (
	"encoding/json"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/segmentio/kafka-go"
)

type Row struct {
	Key    []byte `json:"key,omitempty"`
	Values []any  `json:"values,omitempty"`
}

func GenerateMessages(table Table, num int) (messages []kafka.Message) {
	for i := 0; i < num; i++ {
		messages = append(messages, GenerateMessage(table))
	}
	if len(messages) > 0 {
		log.Printf("generated %d messages values like: %s", num, messages[0].Value)
	}

	return
}

func GenerateMessage(table Table) kafka.Message {
	row := Row{
		Key:    nil,
		Values: nil,
	}

	for i, field := range table.Schema {
		rnd := time.Now().Nanosecond()
		if i == 0 {
			rnd = rnd / 1e6 // limit key range to 0-999
		}
		rndStr := strconv.Itoa(rnd)
		prefix := strings.ToLower(field.Name) + ":"
		var value any
		switch field.Type {
		case FieldTypeString:
			value = prefix + rndStr
			if i == 0 {
				row.Key = []byte(value.(string))
			}
		default:
			panic("TODO: maps all types")
		}
		row.Values = append(row.Values, value)
	}

	data, err := json.Marshal(row)
	if err != nil {
		log.Println("failed to json.Marshal row:", row)
		log.Fatal(err)
	}

	return kafka.Message{
		Key:   row.Key,
		Value: data,
	}
}
