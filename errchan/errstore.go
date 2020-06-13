package errchan

import (
	"encoding/json"
	"strings"
	"sync"
)

type errList []error

type errCollector interface {
	errors() []error
}

type errStore struct {
	collector errCollector
	errList   errList
	mu        sync.Mutex
	collected uint32
}

func newStore(ec errCollector) *errStore {
	return &errStore{ec, nil, sync.Mutex{}, 0}
}

// Errors processes all errors in the collector once and returns the result as slice.
func (s *errStore) Errors() []error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.collected == 0 {
		s.errList = s.collector.errors()
		s.collected = 1
	}
	return s.errList
}

// Strings all errors as []string.
func (s *errStore) Strings() []string {
	errs := s.Errors()
	texts := make([]string, len(errs))
	for i, err := range errs {
		texts[i] = err.Error()
	}
	return texts
}

// String returns errors as strings.
func (s *errStore) String() string {
	return strings.Join(s.Strings(), "\n")
}

// JSON returns errors as JSON list.
func (s *errStore) JSON() []byte {
	errs := s.Strings()
	res, err := json.Marshal(errs)
	if err != nil {
		panic("failed to Marshal error strings" + strings.Join(errs, "\n"))
	}
	return res
}
