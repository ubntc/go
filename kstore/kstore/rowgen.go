package kstore

import (
	"encoding/json"
	"log"
	"strconv"
	"strings"
	"time"
)

func GenerateRows(table *tableSchema, num int) (rows []Row) {
	for i := 0; i < num; i++ {
		row := GenerateRow(table)
		rows = append(rows, row)
	}
	if len(rows) > 0 {
		log.Printf("generated %d rows values like: %v", num, rows[0].Values)
	}
	return
}

func GenerateRow(table *tableSchema) Row {
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
				v, ok := value.(string)
				if !ok {
					panic("value is not a string")
				}
				row.Key = []byte(v)
			}
		default:
			panic("TODO: map more unsupported types")
		}
		row.Values = append(row.Values, value)
	}

	_, err := json.Marshal(row)
	if err != nil {
		log.Println("failed to json.Marshal row:", row)
		log.Fatal(err)
	}

	return row
}
