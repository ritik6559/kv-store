package store

import "errors"

var ErrKeyDoesNotExist = errors.New("key does not exist")
var ErrKeyAlreadyExists = errors.New("key already exists")
var ErrEmptyKey = errors.New("key is mandatory, min length is 1")
var ErrStoreFull = errors.New("store is full")
var ErrInvalidCapacity = errors.New("capacity must be greater than 0")
var ErrInvalidTTL = errors.New("ttl must be greater than 0")

type Store interface {
	Get(key string) (string, error)
	Set(key, value string) error
	Keys() []string
	Delete(key string)
	Len() int
}
