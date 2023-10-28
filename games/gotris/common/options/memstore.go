package options

import "sync"

// MemStore implements options.Options using slices and a map.
// MemStore is fully concurrency-safe and returns copies of
// all value when using any of the getters methods.
type MemStore struct {
	Options       []string
	Descriptions  []string
	values        map[string]string
	currentOption int

	// single changed channel to publish changes
	changed chan bool

	mu sync.RWMutex
}

func NewMemStore(names []string, descs []string) *MemStore {
	return &MemStore{
		Options:      names,
		Descriptions: descs,
		values:       make(map[string]string),
		changed:      make(chan bool, 1),
	}
}

func (s *MemStore) List() (options []string) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return append(options, s.Options...)
}

func (s *MemStore) Descs() (descs []string) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return append(descs, s.Descriptions...)
}

func (s *MemStore) GetName() string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.Options[s.Get()]
}

func (s *MemStore) Values() map[string]string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	cp := make(map[string]string)
	for k, v := range s.values {
		cp[k] = v
	}
	return cp
}

func (s *MemStore) Set(idx int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	n := len(s.Options)
	if n == 0 {
		return
	}

	idx = idx % n       // move into range: [-len:len]
	idx = (idx + n) % n // move into range: [0:len]

	if idx != s.currentOption {
		s.currentOption = idx
		// send change
		select {
		case s.changed <- true:
		default:
			// if the channel is blocked, it means there is a prev. `true` pending
			// no need to send anything then
		}
	}
}

func (s *MemStore) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return len(s.Options)
}

func (s *MemStore) Inc() {
	s.Set(s.currentOption + 1)
}

func (s *MemStore) Dec() {
	s.Set(s.currentOption - 1)
}

func (s *MemStore) Get() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if len(s.Options) == 0 {
		// 0 is the default even if the list is empty
		return 0
	}
	return s.currentOption
}

func (s *MemStore) SetValue(key, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.values[key] = value
}

func (s *MemStore) GetValue(key string) string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.values[key]
}

// Changed returns the stores changed channel.
func (s *MemStore) Changed() <-chan bool {
	return s.changed
}
