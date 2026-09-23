package main

import (
	"encoding/base64"
	"log"

	"github.com/ritik6559/kv-store/internal/store"
)

func main() {
	log.Println("Hello")
}

func SetKeyWithEncryption(store store.Store, key, val string) (string, error) {
	encoded := base64.StdEncoding.EncodeToString([]byte(val))
	if err := store.Set(key, encoded); err != nil {
		return "", err
	}
	return store.Get(key)
}
