package options

// MemStore implements options.Options using slices and a map.
// MemStore is not concurrency-safe.
type MemStore struct {
	Options       []string
	Descriptions  []string
	values        map[string]string
	currentOption int
}

func NewMemStore(names []string, descs []string) *MemStore {
	return &MemStore{
		Options:      names,
		Descriptions: descs,
		values:       make(map[string]string),
	}
}

func (s *MemStore) List() []string {
	return s.Options
}

func (s *MemStore) Descs() []string {
	return s.Descriptions
}

func (s *MemStore) GetName() string {
	return s.Options[s.Get()]
}

func (s *MemStore) Set(idx int) {
	if idx >= len(s.Options) || idx < 0 {
		panic("options index ouf of bounds")
	}
	s.currentOption = idx
}

func (s *MemStore) Get() int {
	if len(s.Options) == 0 {
		panic("cannot access empty options list")
	}
	return s.currentOption
}

func (s *MemStore) Len() int {
	return len(s.Options)
}

func (s *MemStore) Values() map[string]string {
	return s.values
}

func (s *MemStore) SetValue(key, value string) {
	s.values[key] = value
}

func (s *MemStore) GetValue(key string) string {
	return s.values[key]
}
