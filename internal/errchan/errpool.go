package errchan

import (
	"encoding/json"
	"strings"
	"sync"
)

// Pool collects concurrent errors. EXPERIMENTAL! DO NOT USE!
type Pool struct {
	sync.WaitGroup
	sync.Pool
	result []error
}

// NewPool creates a new ChanPool and starts a collector goroutine to collect errors from ChanPool.C.
func NewPool() *Pool {
	return &Pool{}
}

// Wait returns the error channel.
func (c *Pool) Wait() {
	c.WaitGroup.Wait()
	var errs []error
	for {
		v := c.Get()
		if v == nil {
			break
		}
		errs = append(errs, v.(error))
	}
	c.result = errs
}

// Errors returns all errors as slice.
func (c *Pool) Errors() []error {
	return c.result
}

// Strings returns all errors as []string.
func (c *Pool) Strings() []string {
	errs := c.Errors()
	strs := make([]string, len(errs))
	for i, v := range errs {
		strs[i] = v.Error()
	}
	return strs
}

// String returns errors as strings.
func (c *Pool) String() string {
	return strings.Join(c.Strings(), "\n")
}

// JSON returns errors as JSON list.
func (c *Pool) JSON() []byte {
	errs := c.Strings()
	res, err := json.Marshal(errs)
	if err != nil {
		panic("failed to Marshal error strings" + strings.Join(errs, "\n"))
	}
	return res
}
