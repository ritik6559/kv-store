package middleware

import (
	"errors"
	"fmt"

	"github.com/ritik6559/kv-store/internal/store"
)

var _ store.Store = (*MetricsMiddleware)(nil)

type MetricsMiddleware struct {
	inner store.Store

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
	m.getCalls++
	value, err := m.inner.Get(key)
	if errors.Is(err, store.ErrKeyDoesNotExist) {
		m.getMisses++
	}
	return value, err
}

func (m *MetricsMiddleware) GetMisses() int {
	return m.getMisses
}

func (m *MetricsMiddleware) Set(key, value string) error {
	m.setCalls++
	return m.inner.Set(key, value)
}

func (m *MetricsMiddleware) Incr(key string) (int64, error) {
	m.incrCalls++
	return m.inner.Incr(key)
}

func (m *MetricsMiddleware) Keys() []string {
	return m.inner.Keys()
}

func (m *MetricsMiddleware) Delete(key string) {
	m.deleteCalls++
	m.inner.Delete(key)
}

func (m *MetricsMiddleware) Len() int {
	m.lenCalls++
	return m.inner.Len()
}

func (m *MetricsMiddleware) Report() {
	fmt.Printf("Metrics: get_calls=%d get_misses=%d set_calls=%d incr_calls=%d delete_calls=%d len_calls=%d\n",
		m.getCalls,
		m.getMisses,
		m.setCalls,
		m.incrCalls,
		m.deleteCalls,
		m.lenCalls,
	)
}
