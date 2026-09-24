package main

import (
	"encoding/base64"
	"log"

	"github.com/ritik6559/kv-store/internal/store"
)

func main() {
	log.Println("Hello")
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
