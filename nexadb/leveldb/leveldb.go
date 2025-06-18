package leveldb

import (
	"errors"
	"fmt"

	"github.com/syndtr/goleveldb/leveldb"
)

var (
	ErrNotFound = errors.New("item not found in leveldb database")
)

type Database struct {
	db *leveldb.DB
}

func New(path string) (*Database, error) {
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return nil, err
	}
	return &Database{db: db}, nil
}

func (db *Database) Put(key []byte, value []byte) error {
	return db.db.Put(key, value, nil)
}

func (db *Database) Get(key []byte) ([]byte, error) {
	return db.db.Get(key, nil)
}

func (db *Database) Delete(key []byte) error {
	return db.db.Delete(key, nil)
}

func (db *Database) First() ([]byte, []byte, error) {
	iter := db.db.NewIterator(nil, nil)
	defer iter.Release()

	if iter.First() {
		key := append([]byte{}, iter.Key()...)
		value := append([]byte{}, iter.Value()...)
		return key, value, nil
	}
	return nil, nil, ErrNotFound
}

func (db *Database) Last() ([]byte, []byte, error) {
	iter := db.db.NewIterator(nil, nil)
	defer iter.Release()

	if iter.Last() {
		key := append([]byte{}, iter.Key()...)
		value := append([]byte{}, iter.Value()...)
		return key, value, nil
	}
	return nil, nil, ErrNotFound
}

func (db *Database) Previous() ([]byte, []byte, error) {
	iter := db.db.NewIterator(nil, nil)
	defer iter.Release()

	if !iter.Last() {
		return nil, nil, fmt.Errorf("database is empty")
	}

	if !iter.Prev() {
		return nil, nil, fmt.Errorf("only one item in database")
	}

	key := append([]byte{}, iter.Key()...)
	value := append([]byte{}, iter.Value()...)
	return key, value, nil
}
