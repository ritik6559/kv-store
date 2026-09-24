package logging

import (
	"errors"
	"log/slog"
	"time"

	"github.com/ritik6559/kv-store/internal/store"
)

var _ store.Store = (*LoggingStore)(nil)

type LoggingStore struct {
	next   store.Store
	logger *slog.Logger
}

func NewLoggingStore(next store.Store, logger *slog.Logger) *LoggingStore {
	return &LoggingStore{
		next:   next,
		logger: logger,
	}
}

func (l *LoggingStore) Get(key string) (string, error) {
	start := time.Now()
	val, err := l.next.Get(key)
	l.log("get", err, "key", key, "took", time.Since(start))
	return val, err
}

func (l *LoggingStore) Set(key, value string) error {
	start := time.Now()
	err := l.next.Set(key, value)
	l.log("set", err, "key", key, "took", time.Since(start))
	return err
}

func (l *LoggingStore) Keys() []string {
	start := time.Now()
	keys := l.next.Keys()
	l.logger.Debug("keys", "count", len(keys), "took", time.Since(start))
	return keys
}

func (l *LoggingStore) Delete(key string) {
	start := time.Now()
	l.next.Delete(key)
	l.logger.Debug("delete", "key", key, "took", time.Since(start))
}

func (l *LoggingStore) Len() int {
	n := l.next.Len()
	l.logger.Debug("len", "count", n)
	return n
}

func (l *LoggingStore) log(op string, err error, args ...any) {
	switch {
	case err == nil, errors.Is(err, store.ErrKeyDoesNotExist): // OR
		l.logger.Debug(op, append(args, "err", err)...)
	case errors.Is(err, store.ErrEmptyKey), errors.Is(err, store.ErrStoreFull):
		l.logger.Warn(op, append(args, "err", err)...)
	default:
		l.logger.Error(op, append(args, "err", err)...)
	}
}
