package main

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"sync"
)

var store = struct {
	sync.RWMutex
	m map[string]string
}{m: make(map[string]string)}

func Put(key string, value string) error {
	store.Lock()
	defer store.Unlock()

	encodedKey := encodeKey(key)

	store.m[encodedKey] = value

	return nil
}

var ErrorNoSuchKey = errors.New("no such key")

func Get(key string) (string, error) {
	store.RLock()
	defer store.RUnlock()

	encodedKey := encodeKey(key)

	value, ok := store.m[encodedKey]

	if !ok {
		return "", ErrorNoSuchKey
	}

	return value, nil
}

func encodeKey(key string) string {
	keyHash := sha256.Sum256([]byte(key))

	return hex.EncodeToString(keyHash[:])
}

func Delete(key string) error {
	store.Lock()
	defer store.Unlock()

	encodedKey := encodeKey(key)

	delete(store.m, encodedKey)

	return nil
}
