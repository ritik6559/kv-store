package kv

import (
	"fmt"
	"sort"

	"github.com/ritik6559/kv-store/internal/store"
)

var _ store.Store = (*KeyValueStore)(nil)

type KeyValueStore struct {
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
	_, ok := s.data[key]
	if !ok && s.Len() >= s.capacity {
		return fmt.Errorf("%w, capacity is: %d", store.ErrStoreFull, s.capacity)
	}
	s.data[key] = value

	return nil
}

func (s *KeyValueStore) Get(key string) (string, error) {
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
	delete(s.data, key)
}

func (s *KeyValueStore) Len() int {
	return len(s.data)
}

func (s *KeyValueStore) Rename(oldKey, newKey string) error {
	if len(oldKey) == 0 || len(newKey) == 0 {
		return store.ErrEmptyKey
	}
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
	val, ok := s.data[key]
	if !ok {
		return "", store.ErrKeyDoesNotExist
	}

	delete(s.data, key)
	return val, nil
}
