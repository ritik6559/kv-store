package kv

import (
	"fmt"
	"sort"

	"github.com/ritik6559/kv-store/internal/store"
)

type KeyValueStore struct {
	capacity int
	data     map[string]string
}

func NewKeyValueStore(capacity int) *KeyValueStore {
	return &KeyValueStore{
		capacity: capacity,
		data:     make(map[string]string),
	}
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
		return fmt.Errorf("max size reached, capacity is: %d", s.capacity)
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
	val, ok := s.data[oldKey]
	if !ok {
		return store.ErrKeyDoesNotExist
	}

	s.data[newKey] = val
	delete(s.data, oldKey)

	return nil
}

func (s *KeyValueStore) Pop(key string) (string, bool) {
	val, ok := s.data[key]
	if !ok {
		return "", false
	}

	delete(s.data, key)
	return val, true
}
