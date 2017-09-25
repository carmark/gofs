/*
Copyright 2017 The GoStor Authors All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cache

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/boltdb/bolt"
	"github.com/gostor/gofs/pkg/fs"
)

// DB -
type BoltDB struct {
	*bolt.DB
}

func init() {
	if cacheList == nil {
		cacheList = map[string]cacheInitFunc{}
	}
	cacheList["boltdb"] = Open
}

// Open -
func Open(path string, mode os.FileMode) (Cache, error) {
	dname := filepath.Dir(path)
	if err := os.MkdirAll(dname, 0700); err != nil {
		return nil, err
	}
	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		return nil, err
	}

	return &BoltDB{
		db,
	}, nil
}

func (db *BoltDB) Add(ns string, f *fs.File) error {
	tx, err := db.DB.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	bucket := tx.Bucket([]byte(ns))
	if bucket == nil {
		bucket, err = tx.CreateBucket([]byte(ns))
		if err != nil {
			return err
		}
	}
	data, err := json.Marshal(f)
	if err != nil {
		return err
	}
	if err = bucket.Put([]byte(f.FullPath()), data); err != nil {
		return err
	}
	// Commit the transaction and check for error.
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (db *BoltDB) Get(ns, name string) (*fs.File, error) {
	tx, err := db.DB.Begin(true)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	bucket := tx.Bucket([]byte(ns))
	if bucket == nil {
		return nil, fmt.Errorf("No such bucket")
	}
	var f fs.File
	data := bucket.Get([]byte(f.FullPath()))
	if data != nil {
		return nil, fmt.Errorf("No such key")
	}
	err = json.Unmarshal(data, &f)
	if err != nil {
		return nil, err
	}
	// Commit the transaction and check for error.
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return &f, nil
}

func (db *BoltDB) Update(ns, name string, new *fs.File) (*fs.File, error) {
	return nil, nil
}

func (db *BoltDB) Delete(ns, name string) error {
	tx, err := db.DB.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	bucket := tx.Bucket([]byte(ns))
	if bucket == nil {
		return fmt.Errorf("No such bucket")
	}
	if err = bucket.Delete([]byte(name)); err != nil {
		return err
	}
	// Commit the transaction and check for error.
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}
