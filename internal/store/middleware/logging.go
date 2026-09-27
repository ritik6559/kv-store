package middleware

import (
	"errors"
	"log/slog"
	"time"

	"github.com/ritik6559/kv-store/internal/store"
)

var _ store.Store = (*LoggingMiddleware)(nil)

type LoggingMiddleware struct {
	inner  store.Store
	logger *slog.Logger
}

func NewLoggingMiddleware(inner store.Store, logger *slog.Logger) *LoggingMiddleware {
	return &LoggingMiddleware{
		inner:  inner,
		logger: logger,
	}
}

func (l *LoggingMiddleware) Get(key string) (string, error) {
	start := time.Now()
	val, err := l.inner.Get(key)
	l.log("get", err, "key", key, "took", time.Since(start))
	return val, err
}

func (l *LoggingMiddleware) Set(key, value string) error {
	start := time.Now()
	err := l.inner.Set(key, value)
	l.log("set", err, "key", key, "took", time.Since(start))
	return err
}

func (l *LoggingMiddleware) Incr(key string) (int64, error) {
	start := time.Now()
	value, err := l.inner.Incr(key)
	l.log("incr", err, "key", key, "value", value, "took", time.Since(start))
	return value, err
}

func (l *LoggingMiddleware) Keys() []string {
	start := time.Now()
	keys := l.inner.Keys()
	l.logger.Debug("keys", "count", len(keys), "took", time.Since(start))
	return keys
}

func (l *LoggingMiddleware) Delete(key string) {
	start := time.Now()
	l.inner.Delete(key)
	l.logger.Debug("delete", "key", key, "took", time.Since(start))
}

func (l *LoggingMiddleware) Len() int {
	n := l.inner.Len()
	l.logger.Debug("len", "count", n)
	return n
}

func (l *LoggingMiddleware) log(op string, err error, args ...any) {
	switch {
	case err == nil, errors.Is(err, store.ErrKeyDoesNotExist): // OR
		l.logger.Debug(op, append(args, "err", err)...)
	case errors.Is(err, store.ErrEmptyKey), errors.Is(err, store.ErrStoreFull):
		l.logger.Warn(op, append(args, "err", err)...)
	default:
		l.logger.Error(op, append(args, "err", err)...)
	}
}
