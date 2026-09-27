package persistence

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"sync"

	"github.com/ritik6559/kv-store/internal/store"
)

var _ store.Store = (*AOFStore)(nil)

const (
	opSet    = "SET"
	opDelete = "DEL"
	opRename = "RENAME"
)

type record struct {
	Op     string `json:"op"`
	Key    string `json:"key"`
	Value  string `json:"value,omitempty"`
	NewKey string `json:"new_key,omitempty"`
}

type AOFStore struct {
	inner     store.Store
	mu        sync.Mutex
	file      *os.File
	sync      bool
	deleteErr error
}

func Open(path string, inner store.Store, syncEveryWrite bool) (*AOFStore, error) {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0o644)
	if err != nil {
		return nil, fmt.Errorf("opening aof: %w", err)
	}
	err = recoverFile(file, inner)
	if closeErr := file.Close(); err == nil {
		err = closeErr
	}
	if err != nil {
		return nil, err
	}

	file, err = os.OpenFile(path, os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return nil, fmt.Errorf("opening aof: %w", err)
	}

	return &AOFStore{
		inner: inner,
		file:  file,
		sync:  syncEveryWrite,
	}, nil
}

func recoverFile(file *os.File, inner store.Store) error {
	goodEnd, tornTail, err := replay(file, inner)
	if err != nil {
		return err
	}

	info, err := file.Stat()
	if err != nil {
		return fmt.Errorf("stat aof: %w", err)
	}

	switch {
	case tornTail:
		if err := file.Truncate(goodEnd); err != nil {
			return fmt.Errorf("truncating torn aof tail: %w", err)
		}
	case goodEnd > info.Size():
		if _, err := file.WriteAt([]byte{'\n'}, info.Size()); err != nil {
			return fmt.Errorf("repairing aof tail: %w", err)
		}
	}
	return nil
}

func replay(r io.Reader, inner store.Store) (goodEnd int64, tornTail bool, err error) {
	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 64*1024), 16*1024*1024) 

	var offset int64
	lineNo := 0
	var pending error
	for scanner.Scan() {
		lineNo++
		if pending != nil {
			return 0, false, pending
		}

		line := scanner.Bytes()
		offset += int64(len(line)) + 1 

		var rec record
		if err := json.Unmarshal(line, &rec); err != nil {
			pending = fmt.Errorf("aof line %d: %w", lineNo, err)
			continue
		}
		if err := apply(inner, rec); err != nil {
			return 0, false, fmt.Errorf("aof line %d: %w", lineNo, err)
		}
		goodEnd = offset
	}
	if err := scanner.Err(); err != nil {
		return 0, false, fmt.Errorf("reading aof: %w", err)
	}
	return goodEnd, pending != nil, nil
}

func apply(s store.Store, rec record) error {
	switch rec.Op {
	case opSet:
		return s.Set(rec.Key, rec.Value)
	case opDelete:
		s.Delete(rec.Key)
		return nil
	case opRename:
		return s.Rename(rec.Key, rec.NewKey)
	default:
		return fmt.Errorf("unknown op %q", rec.Op)
	}
}

func (a *AOFStore) append(rec record) error {
	line, err := json.Marshal(rec)
	if err != nil {
		return err
	}
	line = append(line, '\n')

	if _, err := a.file.Write(line); err != nil {
		return fmt.Errorf("writing aof: %w", err)
	}
	if a.sync {
		if err := a.file.Sync(); err != nil {
			return fmt.Errorf("syncing aof: %w", err)
		}
	}
	return nil
}

func (a *AOFStore) Set(key, value string) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if err := a.inner.Set(key, value); err != nil {
		return err
	}
	return a.append(record{Op: opSet, Key: key, Value: value})
}

func (a *AOFStore) Incr(key string) (int64, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	n, err := a.inner.Incr(key)
	if err != nil {
		return 0, err
	}
	return n, a.append(record{Op: opSet, Key: key, Value: strconv.FormatInt(n, 10)})
}

func (a *AOFStore) Delete(key string) bool {
	a.mu.Lock()
	defer a.mu.Unlock()

	existed := a.inner.Delete(key)
	if existed {
		if err := a.append(record{Op: opDelete, Key: key}); err != nil {
			a.deleteErr = errors.Join(a.deleteErr, err)
		}
	}
	return existed
}

func (a *AOFStore) Rename(oldKey, newKey string) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if err := a.inner.Rename(oldKey, newKey); err != nil {
		return err
	}
	if oldKey == newKey {
		return nil
	}
	return a.append(record{Op: opRename, Key: oldKey, NewKey: newKey})
}

func (a *AOFStore) Pop(key string) (string, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	val, err := a.inner.Pop(key)
	if err != nil {
		return "", err
	}
	return val, a.append(record{Op: opDelete, Key: key})
}

func (a *AOFStore) Get(key string) (string, error) { return a.inner.Get(key) }
func (a *AOFStore) Keys() []string                 { return a.inner.Keys() }
func (a *AOFStore) Len() int                       { return a.inner.Len() }

func (a *AOFStore) Close() error {
	a.mu.Lock()
	defer a.mu.Unlock()

	return errors.Join(a.deleteErr, a.file.Sync(), a.file.Close())
}
