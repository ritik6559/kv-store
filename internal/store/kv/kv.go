package kv

import (
	"fmt"
	"math"
	"sort"
	"strconv"
	"sync"

	"github.com/ritik6559/kv-store/internal/store"
)

var _ store.Store = (*KeyValueStore)(nil)

type KeyValueStore struct {
	mu       sync.RWMutex
	capacity int
	data     map[string]string
}

func NewKeyValueStore(capacity int) (*KeyValueStore, error) {
	if capacity <= 0 {
		return nil, store.ErrInvalidCapacity
	}
	return &KeyValueStore{
		capacity: capacity,
		data:     make(map[string]string),
	}, nil
}

func (s *KeyValueStore) Keys() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	keys := make([]string, 0, len(s.data))

	for key := range s.data {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	return keys
}

func (s *KeyValueStore) Set(key, value string) error {
	if len(key) == 0 {
		return store.ErrEmptyKey
	}
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.data[key]
	if !ok && len(s.data) >= s.capacity {
		return fmt.Errorf("%w, capacity is: %d", store.ErrStoreFull, s.capacity)
	}
	s.data[key] = value

	return nil
}

func (s *KeyValueStore) Incr(key string) (int64, error) {
	if len(key) == 0 {
		return 0, store.ErrEmptyKey
	}
	s.mu.Lock()
	defer s.mu.Unlock()

	current := int64(0)
	value, exists := s.data[key]
	if exists {
		parsed, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("value for key %q is not an integer: %w", key, err)
		}
		current = parsed
	} else if len(s.data) >= s.capacity {
		return 0, fmt.Errorf("%w, capacity is: %d", store.ErrStoreFull, s.capacity)
	}
	if current == math.MaxInt64 {
		return 0, fmt.Errorf("incrementing key %q: integer overflow", key)
	}

	current++
	s.data[key] = strconv.FormatInt(current, 10)
	return current, nil
}

func (s *KeyValueStore) Get(key string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if len(key) == 0 {
		return "", store.ErrEmptyKey
	}
	val, ok := s.data[key]
	if !ok {
		return "", store.ErrKeyDoesNotExist
	}
	return val, nil
}

func (s *KeyValueStore) Delete(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.data, key)
}

func (s *KeyValueStore) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return len(s.data)
}

func (s *KeyValueStore) Rename(oldKey, newKey string) error {
	if len(oldKey) == 0 || len(newKey) == 0 {
		return store.ErrEmptyKey
	}
	s.mu.Lock()
	defer s.mu.Unlock()

	val, ok := s.data[oldKey]
	if !ok {
		return store.ErrKeyDoesNotExist
	}
	if oldKey == newKey {
		return nil
	}
	if _, exists := s.data[newKey]; exists {
		return store.ErrKeyAlreadyExists
	}

	s.data[newKey] = val
	delete(s.data, oldKey)

	return nil
}

func (s *KeyValueStore) Pop(key string) (string, error) {
	if len(key) == 0 {
		return "", store.ErrEmptyKey
	}
	s.mu.Lock()
	defer s.mu.Unlock()

	val, ok := s.data[key]
	if !ok {
		return "", store.ErrKeyDoesNotExist
	}

	delete(s.data, key)
	return val, nil
}
