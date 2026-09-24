package main

import (
	"encoding/base64"
	"log"
	"log/slog"
	"os"

	"github.com/ritik6559/kv-store/internal/store"
	"github.com/ritik6559/kv-store/internal/store/kv"
	"github.com/ritik6559/kv-store/internal/store/logging"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	base, err := kv.NewKeyValueStore(100)
	if err != nil {
		log.Fatal(err)
	}
	s := logging.NewLoggingStore(base, logger.With("store", "kv"))

	if err := SetKeyEncoded(s, "name", "ritik"); err != nil {
		logger.Error("set failed", "err", err)
	}
	if val, err := GetKeyDecoded(s, "name"); err == nil {
		logger.Info("decoded value", "key", "name", "value", val)
	}
	_, _ = s.Get("missing")
	s.Delete("name")
	s.Len()
}

func SetKeyEncoded(s store.Store, key, val string) error {
	encoded := base64.StdEncoding.EncodeToString([]byte(val))
	return s.Set(key, encoded)
}

func GetKeyDecoded(s store.Store, key string) (string, error) {
	encoded, err := s.Get(key)
	if err != nil {
		return "", err
	}
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", err
	}
	return string(decoded), nil
}
