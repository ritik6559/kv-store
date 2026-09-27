package middleware

import (
	"errors"
	"fmt"
	"sync"

	"github.com/ritik6559/kv-store/internal/store"
)

var _ store.Store = (*MetricsMiddleware)(nil)

type MetricsMiddleware struct {
	inner store.Store
	mu    sync.Mutex

	getCalls    int
	getMisses   int
	setCalls    int
	incrCalls   int
	deleteCalls int
	lenCalls    int
}

func NewMetricsMiddleware(inner store.Store) *MetricsMiddleware {
	return &MetricsMiddleware{inner: inner}
}

func (m *MetricsMiddleware) Get(key string) (string, error) {
	m.mu.Lock()
	m.getCalls++
	m.mu.Unlock()

	value, err := m.inner.Get(key)
	if errors.Is(err, store.ErrKeyDoesNotExist) {
		m.mu.Lock()
		m.getMisses++
		m.mu.Unlock()
	}
	return value, err
}

func (m *MetricsMiddleware) GetMisses() int {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.getMisses
}

func (m *MetricsMiddleware) Set(key, value string) error {
	m.mu.Lock()
	m.setCalls++
	m.mu.Unlock()

	return m.inner.Set(key, value)
}

func (m *MetricsMiddleware) Incr(key string) (int64, error) {
	m.mu.Lock()
	m.incrCalls++
	m.mu.Unlock()

	return m.inner.Incr(key)
}

func (m *MetricsMiddleware) Keys() []string {
	return m.inner.Keys()
}

func (m *MetricsMiddleware) Delete(key string) {
	m.mu.Lock()
	m.deleteCalls++
	m.mu.Unlock()

	m.inner.Delete(key)
}

func (m *MetricsMiddleware) Len() int {
	m.mu.Lock()
	m.lenCalls++
	m.mu.Unlock()

	return m.inner.Len()
}

func (m *MetricsMiddleware) Report() {
	m.mu.Lock()
	getCalls := m.getCalls
	getMisses := m.getMisses
	setCalls := m.setCalls
	incrCalls := m.incrCalls
	deleteCalls := m.deleteCalls
	lenCalls := m.lenCalls
	m.mu.Unlock()

	fmt.Printf("Metrics: get_calls=%d get_misses=%d set_calls=%d incr_calls=%d delete_calls=%d len_calls=%d\n",
		getCalls,
		getMisses,
		setCalls,
		incrCalls,
		deleteCalls,
		lenCalls,
	)
}
