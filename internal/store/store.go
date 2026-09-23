package store

import "errors"

var ErrKeyDoesNotExist = errors.New("key does not exist")
var ErrEmptyKey = errors.New("key is mandatory, min length is 1")

type Store interface {
	Get(key string) (string, error)
	Set(key, value string ) error
	Keys() []string
	Delete(key string)
	Len() int
}