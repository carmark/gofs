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

// Package master contains the main function.
package master

import (
	"github.com/gostor/gofs/pkg/cache"
	"github.com/gostor/gofs/pkg/fs"
	"github.com/gostor/gofs/pkg/raft"
)

type Master struct {
	Name       string
	RaftServer *raft.RaftServer
	Namespaces map[string]*fs.Namespace
	Cache      cache.Cache
}

func NewMaster(cfg *MasterConfig) (*Master, error) {
	cc, err := cache.NewCache(cfg.CacheType, cfg.CacheDir, 0700)
	if err != nil {
		return nil, err
	}
	rs, err := raft.NewRaftServer(cfg.Peers, cfg.HttpAddr, cfg.DataDir, cfg.PuleSeconds, cfg.Router, cfg.HttpServers[0], cc)
	if err != nil {
		return nil, err
	}
	return &Master{
		Name:       cfg.HttpAddr,
		RaftServer: rs,
		Namespaces: map[string]*fs.Namespace{},
		Cache:      cc,
	}, nil
}

func (m *Master) IsLeader() bool {
	if l, err := m.RaftServer.Leader(); err == nil {
		return l == m.Name
	}
	return false
}

func (m *Master) GetPathHandler(path, op string) (interface{}, error) {
	return nil
}

func (m *Master) PutPathHandler() error {
	return nil
}

func (m *Master) PostPathHandler(path, op string) error {
	return nil
}

func (m *Master) DeletePathHandler(path, op string, r bool) error {
	return nil
}
