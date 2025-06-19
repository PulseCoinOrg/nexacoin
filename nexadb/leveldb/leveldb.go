/*
 * NexaCoin - A Cryptocurrency Framework
 *
 * Copyright (c) 2025 NexaCoin Developers
 *
 * This file is part of the NexaCoin project and is licensed under the MIT License.
 * You may obtain a copy of the License at:
 *
 *     https://opensource.org/licenses/MIT
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
 * THE SOFTWARE.
 */

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

// inserts bytes into the leveldb database
func (db *Database) Put(key []byte, value []byte) error {
	return db.db.Put(key, value, nil)
}

// retrieves a value from the database given a key
func (db *Database) Get(key []byte) ([]byte, error) {
	return db.db.Get(key, nil)
}

// removes an item from the database given a key
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
