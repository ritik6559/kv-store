package middleware

import (
	"errors"
	"sync/atomic"

	"github.com/ritik6559/kv-store/internal/store"
)

var _ store.Store = (*MetricsMiddleware)(nil)

type MetricsMiddleware struct {
	inner store.Store

	getCalls    atomic.Int64
	getMisses   atomic.Int64
	setCalls    atomic.Int64
	incrCalls   atomic.Int64
	keysCalls   atomic.Int64
	deleteCalls atomic.Int64
	lenCalls    atomic.Int64
	renameCalls atomic.Int64
	popCalls    atomic.Int64
}

type Stats struct {
	GetCalls    int64
	GetMisses   int64
	SetCalls    int64
	IncrCalls   int64
	KeysCalls   int64
	DeleteCalls int64
	LenCalls    int64
	RenameCalls int64
	PopCalls    int64
}

func NewMetricsMiddleware(inner store.Store) *MetricsMiddleware {
	return &MetricsMiddleware{inner: inner}
}

func (m *MetricsMiddleware) Get(key string) (string, error) {
	m.getCalls.Add(1)
	value, err := m.inner.Get(key)
	if errors.Is(err, store.ErrKeyDoesNotExist) {
		m.getMisses.Add(1)
	}
	return value, err
}

func (m *MetricsMiddleware) Set(key, value string) error {
	m.setCalls.Add(1)
	return m.inner.Set(key, value)
}

func (m *MetricsMiddleware) Incr(key string) (int64, error) {
	m.incrCalls.Add(1)
	return m.inner.Incr(key)
}

func (m *MetricsMiddleware) Keys() []string {
	m.keysCalls.Add(1)
	return m.inner.Keys()
}

func (m *MetricsMiddleware) Delete(key string) bool {
	m.deleteCalls.Add(1)
	return m.inner.Delete(key)
}

func (m *MetricsMiddleware) Len() int {
	m.lenCalls.Add(1)
	return m.inner.Len()
}

func (m *MetricsMiddleware) Rename(oldKey, newKey string) error {
	m.renameCalls.Add(1)
	return m.inner.Rename(oldKey, newKey)
}

func (m *MetricsMiddleware) Pop(key string) (string, error) {
	m.popCalls.Add(1)
	return m.inner.Pop(key)
}

func (m *MetricsMiddleware) Stats() Stats {
	return Stats{
		GetCalls:    m.getCalls.Load(),
		GetMisses:   m.getMisses.Load(),
		SetCalls:    m.setCalls.Load(),
		IncrCalls:   m.incrCalls.Load(),
		KeysCalls:   m.keysCalls.Load(),
		DeleteCalls: m.deleteCalls.Load(),
		LenCalls:    m.lenCalls.Load(),
		RenameCalls: m.renameCalls.Load(),
		PopCalls:    m.popCalls.Load(),
	}
}
