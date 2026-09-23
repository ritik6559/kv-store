package main

import (
	"fmt"
	"time"
)

type TTLStore struct {
	data map[string]ttlEntry
}

type ttlEntry struct {
	value     string
	expiresAt time.Time
}

func NewTTLStore() *TTLStore {
	return &TTLStore{
		data: make(map[string]ttlEntry),
	}
}

func (t *TTLStore) Set(key, value string, ttl time.Duration) error {
	if key == "" {
		return ErrEmptyKey
	}

	return nil
}

func (t *TTLStore) GetWithTTL(key string) (string, error) {
	if key == "" {
		return "", ErrEmptyKey
	}

	entry, ok := t.data[key]
	if !ok || time.Now().After(entry.expiresAt) {
		// laxy deletion
		delete(t.data, key)
		return "", fmt.Errorf("key expired")
	}

	return entry.value, nil
}

func (t *TTLStore) Delete(key string) {
	delete(t.data, key)
}

func (t *TTLStore) Len() int {
	return len(t.data)
}
