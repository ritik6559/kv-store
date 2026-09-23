package ttl

import (
	"fmt"
	"time"

	"github.com/ritik6559/kv-store/internal/store"
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
		return store.ErrEmptyKey
	}

	entry := ttlEntry{
		value:     value,
		expiresAt: time.Now().Add(ttl),
	}
	t.data[key] = entry

	return nil
}

func (t *TTLStore) Get(key string) (string, error) {
	if key == "" {
		return "", store.ErrEmptyKey
	}

	entry, ok := t.data[key]
	if !ok || time.Now().After(entry.expiresAt) {
		// laxy deletion
		delete(t.data, key)
		return "", fmt.Errorf("key does not exists")
	}

	return entry.value, nil
}

func (t *TTLStore) Delete(key string) {
	delete(t.data, key)
}

func (t *TTLStore) Len() int {
	return len(t.data)
}
