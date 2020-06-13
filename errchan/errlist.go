package errchan

import (
	"sync"
)

// ErrList collects concurrent errors.
type ErrList struct {
	errs []error
	sync.WaitGroup
	sync.Mutex
	// lock   func()
	// unlock func()
}

// NewList creates a new errchan.List.
func NewList() *ErrList { return &ErrList{} }

// Append add an error.
func (c *ErrList) Append(err error) {
	c.Lock()
	defer c.Unlock()
	c.errs = append(c.errs, err)
}
