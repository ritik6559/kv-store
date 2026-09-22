package main

import (
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