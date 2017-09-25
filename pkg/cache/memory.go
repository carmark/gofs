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
	"fmt"
	"os"
	"sync"

	"github.com/gostor/gofs/pkg/fs"
)

// DB -
type Memory struct {
	Namespace map[string]*fs.Namespace
	Files     map[string]map[string]*fs.File
	lock      sync.RWMutex
}

func init() {
	if cacheList == nil {
		cacheList = map[string]cacheInitFunc{}
	}
	cacheList["memory"] = memoryOpen
}

// Open -
func memoryOpen(path string, mode os.FileMode) (Cache, error) {
	return &Memory{
		Namespace: map[string]*fs.Namespace{},
		Files:     map[string]map[string]*fs.File{},
	}, nil
}

func (db *Memory) Add(ns string, f *fs.File) error {
	db.lock.RLock()
	defer db.lock.RUnlock()
	if _, ok := db.Files[ns]; !ok {
		db.Files[ns] = map[string]*fs.File{}
	}
	db.Files[ns][f.FullPath()] = f
	return nil
}

func (db *Memory) Get(ns, name string) (*fs.File, error) {
	var f *fs.File
	db.lock.RLock()
	defer db.lock.RUnlock()
	if files, ok := db.Files[ns]; ok {
		f = files[name]
	} else {
		return nil, fmt.Errorf("no such file")
	}
	return f, nil
}

func (db *Memory) Update(ns, name string, new *fs.File) (*fs.File, error) {
	db.lock.Lock()
	defer db.lock.Unlock()
	if files, ok := db.Files[ns]; !ok {
		return nil, fmt.Errorf("no such file")
	} else {
		if _, ok = files[name]; !ok {
			return nil, fmt.Errorf("no such file")
		}
	}
	db.Files[ns][new.FullPath()] = new
	return nil, nil
}

func (db *Memory) Delete(ns, name string) error {
	db.lock.Lock()
	defer db.lock.Unlock()
	if files, ok := db.Files[ns]; !ok {
		return fmt.Errorf("no such file")
	} else {
		if _, ok = files[name]; !ok {
			return fmt.Errorf("no such file")
		}
	}
	delete(db.Files[ns], name)
	return nil
}
