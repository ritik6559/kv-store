package ttl

import (
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

func (e ttlEntry) expired(now time.Time) bool {
	return now.After(e.expiresAt)
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
	if ttl <= 0 {
		return store.ErrInvalidTTL
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
	if !ok {
		return "", store.ErrKeyDoesNotExist
	}
	if entry.expired(time.Now()) {
		// lazy deletion
		delete(t.data, key)
		return "", store.ErrKeyDoesNotExist
	}

	return entry.value, nil
}

func (t *TTLStore) Delete(key string) {
	delete(t.data, key)
}

// Len returns the number of live keys, removing any expired ones it finds.
func (t *TTLStore) Len() int {
	now := time.Now()
	for key, entry := range t.data {
		if entry.expired(now) {
			delete(t.data, key)
		}
	}
	return len(t.data)
}
