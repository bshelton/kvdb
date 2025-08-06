package storage

// Storage handles the core key-value storage and value counting
type Storage struct {
	data        map[string]string
	valueCounts map[string]int
}

// New creates a new storage instance
func New() *Storage {
	return &Storage{
		data:        make(map[string]string),
		valueCounts: make(map[string]int),
	}
}

// Set stores a key-value pair
func (s *Storage) Set(key, value string) {
	s.data[key] = value
}

// Get retrieves a value by key, returns "NULL" if not found
func (s *Storage) Get(key string) string {
	if value, exists := s.data[key]; exists {
		return value
	}
	return "NULL"
}

// Unset removes a key-value pair
func (s *Storage) Unset(key string) {
	delete(s.data, key)
}

// GetValueCount returns the count of keys with the given value
func (s *Storage) GetValueCount(value string) int {
	return s.valueCounts[value]
}

// UpdateValueCount handles value count changes when a key is updated
func (s *Storage) UpdateValueCount(oldValue, newValue string) {
	if oldValue != "NULL" {
		s.DecrementValueCount(oldValue)
	}
	s.IncrementValueCount(newValue)
}

// IncrementValueCount increases the count for a value
func (s *Storage) IncrementValueCount(value string) {
	s.valueCounts[value]++
}

// DecrementValueCount decreases the count for a value
func (s *Storage) DecrementValueCount(value string) {
	s.valueCounts[value]--
	if s.valueCounts[value] <= 0 {
		delete(s.valueCounts, value)
	}
}