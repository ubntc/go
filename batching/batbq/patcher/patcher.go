package patcher

import (
	"context"
	"encoding/json"
	"log"

	"cloud.google.com/go/bigquery"
)

// FieldMap stores fields.
type FieldMap map[string]*Field

// Field stores field infos.
type Field struct {
	fieldSchema *bigquery.FieldSchema
	fields      FieldMap
}

func fieldMap(schema bigquery.Schema) FieldMap {
	res := make(FieldMap)
	for _, f := range schema {
		field := &Field{fieldSchema: f}
		if f.Type == bigquery.RecordFieldType {
			field.fields = fieldMap(f.Schema)
		}
		res[f.Name] = field
	}
	return res
}

func fieldSchema(m FieldMap) bigquery.Schema {
	res := make(bigquery.Schema, 0, len(m))
	for _, v := range m {
		f := v.fieldSchema
		if f.Type == bigquery.RecordFieldType {
			f.Schema = fieldSchema(v.fields)
		}
		// cannot add required fields to existing schema
		f.Required = false
		res = append(res, f)
	}
	return res
}

func copyFields(src, trg map[string]*Field) bool {
	updated := false
	for name, field := range src {
		if _, ok := trg[name]; !ok {
			trg[name] = field
			updated = true
			continue
		}
		if field.fieldSchema.Type == bigquery.RecordFieldType {
			if copyFields(field.fields, trg[name].fields) {
				updated = true
			}
		}
	}
	return updated
}

func mergeSchema(a, b bigquery.Schema) (bigquery.Schema, bool) {
	src := fieldMap(a)
	trg := fieldMap(b)
	updated := copyFields(src, trg)
	return fieldSchema(trg), updated
}

// Patch patches a table.
func Patch(ctx context.Context, table *bigquery.Table, schema bigquery.Schema) error {
	meta, err := table.Metadata(ctx)
	if err != nil {
		return err
	}
	newSchema, updated := mergeSchema(schema, meta.Schema)
	data, err := json.Marshal(newSchema)
	if err != nil {
		return err
	}

	if !updated {
		log.Printf("schema did not change: schema=%s", string(data))
		return nil
	}
	log.Printf("patching table %s", table.TableID)
	_, err = table.Update(ctx, bigquery.TableMetadataToUpdate{Schema: newSchema}, "")
	return err
}
