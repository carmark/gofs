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
package raft

import (
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/goraft/raft"
	"github.com/gostor/gofs/pkg/cache"
	"github.com/gostor/gofs/pkg/fs"
)

// This command writes a value to a key.
type Operation struct {
	Type      string    `json:"type"`
	Namespace string    `json:"namespace"`
	Filename  string    `json:"name"`
	NewName   string    `json:"newname"`
	FileAttr  *fs.Attr  `json:"attr"`
	CreatedAt time.Time `json:"createdat"`
}

// Creates a new operation command.
func NewOperation(t, n, f, nn string, attr *fs.Attr, created time.Time) *Operation {
	return &Operation{
		Type:      t,
		Namespace: n,
		Filename:  f,
		NewName:   nn,
		FileAttr:  attr,
		CreatedAt: created,
	}
}

// The name of the command in the log.
func (o *Operation) CommandName() string {
	return "metadata_update_operation"
}

// Operate the file of the namespace.
func (o *Operation) Apply(server raft.Server) (interface{}, error) {
	log.Debugf("Raft Apply: [Type: %v, Namespace: %v, Filename: %v, Attr: [%#v]]", o.Type, o.Namespace, o.Filename, o.FileAttr)
	cache := server.Context().(cache.Cache)
	_ = cache
	return nil, nil
}
