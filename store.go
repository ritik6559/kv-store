package main

import (
	"fmt"
	"sort"
)

type store struct {
	data map[string]string
}

func NewStore() *store {
	return &store{
		data: make(map[string]string),
	}
}

func (s *store) Keys() []string {
	keys := make([]string, 0, len(s.data))

	for key := range s.data {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	return keys
}

func (s *store) Set(key, value string) {
	s.data[key] = value
}

func (s *store) Get(key string) (string, bool) {
	val, ok := s.data[key]
	return val, ok
}

func (s *store) Delete(key string) {
	delete(s.data, key)
}

func (s *store) Len() int {
	return len(s.data)
}

func (s *store) Rename(oldKey, newKey string) error {
	val, ok := s.data[oldKey]
	if !ok {
		return fmt.Errorf("key doesn't exists")
	}
	
	s.data[newKey] = val
	delete(s.data, oldKey)

	return nil
}

func (s *store) Pop(key string) (string, bool) {
	val, ok := s.data[key]
	if !ok {
		return "", false
	}

	delete(s.data, key)
	return val, true
}