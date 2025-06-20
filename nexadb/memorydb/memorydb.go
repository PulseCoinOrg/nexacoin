package memorydb

import (
	"errors"
	"sync"
)

var (
	errMemorydbClosed   = errors.New("database closed")
	errMemorydbNotFound = errors.New("not found")
)

type Database struct {
	items map[string][]byte
	lock  sync.RWMutex
}

func New() *Database {
	return &Database{
		items: make(map[string][]byte),
	}
}

func (db *Database) Close() error {
	db.lock.Lock()
	defer db.lock.Unlock()

	db.items = nil
	return nil
}

// Has retrieves if a key is present in the key-value store.
func (db *Database) Has(key []byte) (bool, error) {
	db.lock.RLock()
	defer db.lock.RUnlock()

	if db.items == nil {
		return false, errMemorydbClosed
	}
	_, ok := db.items[string(key)]
	return ok, nil
}

// Get retrieves the given key if it's present in the key-value store.
func (db *Database) Get(key []byte) ([]byte, error) {
	db.lock.RLock()
	defer db.lock.RUnlock()

	if db.items == nil {
		return nil, errMemorydbClosed
	}
	if _, ok := db.items[string(key)]; ok {
		return db.items[string(key)], nil
	}
	return nil, errMemorydbNotFound
}

// Put inserts the given value into the key-value store.
func (db *Database) Put(key []byte, value []byte) error {
	db.lock.Lock()
	defer db.lock.Unlock()

	if db.items == nil {
		return errMemorydbClosed
	}
	db.items[string(key)] = value
	return nil
}

// Delete removes the key from the key-value store.
func (db *Database) Delete(key []byte) error {
	db.lock.Lock()
	defer db.lock.Unlock()

	if db.items == nil {
		return errMemorydbClosed
	}
	delete(db.items, string(key))
	return nil
}
