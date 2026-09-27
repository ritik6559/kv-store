# kv-store

An in-memory key-value store written in Go, with append-only-file persistence, logging, and metrics built as composable middleware.

This is a learning project. It focuses on clean package layout, interfaces, concurrency safety, and crash-safe persistence using only the standard library.

## Features

- **Thread-safe KV store**: `sync.RWMutex`-guarded map with a fixed capacity.
- **Commands**: `GET`, `SET`, `INCR`, `KEYS`, `DELETE`, `LEN`, `RENAME`, `POP`.
- **Persistence (AOF)**: every successful write is appended to `data.aof` as a JSON line and replayed on startup.
  - Recovers from a crash mid-write by truncating a torn last line.
  - Optional fsync after every write.
- **Middleware**: logging (`log/slog`) and metrics (`sync/atomic` counters) wrap any `store.Store`.
- **TTL store**: a separate store with per-key expiry and lazy deletion.
- **Typed errors**: sentinel errors (`ErrKeyDoesNotExist`, `ErrStoreFull`, ...) that work with `errors.Is`.

## Architecture

Every layer implements the same `store.Store` interface, so layers stack like HTTP middleware:

```
commands ──> Logging ──> Metrics ──> AOF persistence ──> KeyValueStore (map)
                                            │
                                            └──> data.aof
```

```go
type Store interface {
	Get(key string) (string, error)
	Set(key, value string) error
	Incr(key string) (int64, error)
	Keys() []string
	Delete(key string) bool
	Len() int
	Rename(oldKey, newKey string) error
	Pop(key string) (string, error)
}
```

Each layer does one job and doesn't know about the others. Any layer can be added, removed, or reordered in `main.go` without touching the store itself.

## Project layout

```
cmd/kv-store/
  main.go             wires the layers together and runs a demo
  command.go          command dispatch (GET, SET, INCR, ...)
internal/store/
  store.go            Store interface and sentinel errors
  kv/kv.go            in-memory store with capacity limit
  ttl/ttl.go          store with per-key expiry
  persistence/aof.go  append-only-file persistence and replay
  middleware/
    logging.go        structured logging with slog
    metrics.go        per-operation call counters
```

## Getting started

Requires Go 1.26+.

```bash
git clone https://github.com/ritik6559/kv-store.git
```

```bash
go run ./cmd/kv-store
```

The demo runs a fixed list of commands, logs each one, and prints metrics at the end. Data is written to `data.aof` in the working directory.

## Commands

| Command  | Example                 | Result                                                  |
| -------- | ----------------------- | ------------------------------------------------------- |
| `SET`    | `SET name ritik`        | Stores a value. Fails if the store is full.             |
| `GET`    | `GET name`              | Returns the value, or `key does not exist`.             |
| `INCR`   | `INCR counter`          | Increments an integer value, starting at 0 if missing.  |
| `DELETE` | `DELETE name`           | Prints `1` if the key existed, `0` otherwise.           |
| `RENAME` | `RENAME name user`      | Renames a key. Fails if the new key already exists.     |
| `POP`    | `POP user`              | Returns the value and deletes the key.                  |
| `KEYS`   | `KEYS`                  | Lists all keys, sorted.                                 |
| `LEN`    | `LEN`                   | Number of keys.                                         |

## Persistence

Each write is stored as one JSON line:

```json
{"op":"SET","key":"name","value":"ritik"}
{"op":"SET","key":"counter","value":"2"}
{"op":"RENAME","key":"name","new_key":"user"}
{"op":"DEL","key":"user"}
```

Design notes:

- **Only state changes are logged.** Reads and failed writes are not.
- **Records are idempotent.** `INCR` is stored as a `SET` of the result and `POP` as a `DEL`, so replay never recomputes anything.
- **Order is preserved.** A mutex covers both the in-memory write and the file append, so the log order always matches the apply order.
- **Crash recovery.** A torn final line (a crash mid-write) is truncated on startup. A bad line anywhere else is treated as corruption and startup fails.


