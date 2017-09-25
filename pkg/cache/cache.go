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

// Package cache maintains the caching of all meta data of the files and directories.
package cache

import (
	"fmt"
	"os"

	"github.com/gostor/gofs/pkg/fs"
)

type Cache interface {
	Add(ns string, f *fs.File) error
	Get(ns, name string) (*fs.File, error)
	Update(ns, name string, new *fs.File) (*fs.File, error)
	Delete(ns, name string) error
}

type cacheInitFunc func(p string, m os.FileMode) (Cache, error)

var cacheList map[string]cacheInitFunc

func NewCache(t, p string, mode os.FileMode) (Cache, error) {
	if fun, ok := cacheList[t]; ok {
		return fun(p, mode)
	}

	return nil, fmt.Errorf("Unknown cache type: %v", t)
}
