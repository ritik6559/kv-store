package main

import (
	"encoding/base64"
	"errors"
	"fmt"
	"sort"
)

var ErrKeyDoesNotExist = errors.New("key does not exist")
var ErrEmptyKey = errors.New("key is mandatory, min length is 1")

type store struct {
	capacity int
	data     map[string]string
}

func NewStore(capacity int) *store {
	return &store{
		capacity: capacity,
		data:     make(map[string]string),
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

func (s *store) SetKeyWithEncryption(key, val string) (string, error) {
	encoded := base64.StdEncoding.EncodeToString([]byte(val))
	if err := s.Set(key, encoded); err != nil {
		return "", err
	}
	return s.Get(key)
}

func (s *store) Set(key, value string) error {
	if len(key) == 0 {
		return ErrEmptyKey
	}
	_, ok := s.data[key]
	if !ok && s.Len() >= s.capacity {
		return fmt.Errorf("max size reached, capacity is: %d", s.capacity)
	}
	s.data[key] = value

	return nil
}

func (s *store) Get(key string) (string, error) {
	if len(key) == 0 {
		return "", ErrEmptyKey
	}
	val, ok := s.data[key]
	if !ok {
		return "", ErrKeyDoesNotExist
	}
	return val, nil
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
		return ErrKeyDoesNotExist
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
