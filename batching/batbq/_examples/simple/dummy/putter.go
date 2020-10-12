package dummy

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"cloud.google.com/go/bigquery"
)

// Putter is a dummy putter.
type Putter struct {
	Name string
	sync.Mutex
	Table        []map[string]bigquery.Value
	WriteDelay   time.Duration
	NumBatches   int
	MaxWorkers   int
	NumErrors    int
	InsertErrors []error
	FatalErr     error
}

// Put stores dummy data.
func (p *Putter) Put(ctx context.Context, src interface{}) error {
	rows := src.([]bigquery.ValueSaver)
	time.Sleep(p.WriteDelay)
	p.Lock()
	defer p.Unlock()
	errors := make(bigquery.PutMultiError, 0)
	for i, v := range rows {
		row, insertID, err := v.Save()
		if err == nil {
			// simulate expected errors
			if strings.HasPrefix(insertID, "fatal") {
				p.FatalErr = fmt.Errorf(insertID)
				return p.FatalErr
			}
			if strings.HasPrefix(insertID, "err") {
				err = fmt.Errorf(insertID)
			}
		}
		if err != nil {
			errors = append(errors, bigquery.RowInsertionError{
				RowIndex: i,
				InsertID: insertID,
				Errors:   bigquery.MultiError{err},
			})
			continue
		}
		p.Table = append(p.Table, row)
	}
	p.NumBatches++
	if len(errors) > 0 {
		p.InsertErrors = append(p.InsertErrors, errors)
		return errors
	}
	return nil
}

// GetNumBatches returns the number of stored batches.
func (p *Putter) GetNumBatches() int {
	p.Lock()
	defer p.Unlock()
	return p.NumBatches
}

// GetLength returns the table size.
func (p *Putter) GetLength() int {
	p.Lock()
	defer p.Unlock()
	return len(p.Table)
}
