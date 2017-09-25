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

package fs

import (
	"context"
	"path"
	"path/filepath"

	"github.com/gostor/gofs/pkg/api"
)

// File implements both Node and Handle for the hello file.
type File struct {
	Attr

	Parent    *File
	Directory bool
	Symlink   bool
	Link      bool
	Path      string
	Checksum  string
	Hash      []byte

	namespace *Namespace
}

func (f *File) IsDirectory() bool {
	return f.Directory
}

func (f *File) IsSymlink() bool {
	return f.Symlink
}

func (f *File) IsLink() bool {
	return f.Link
}

// Setattr - set attribute.
func (f *File) Setattr(ctx context.Context, req *api.SetattrRequest) error {
	// update cache with new attributes
	if req.Valid.Mode() {
		f.Mode = req.Mode
	}

	if req.Valid.Uid() {
		//f.UID = req.Uid
	}

	if req.Valid.Gid() {
		//f.GID = req.Gid
	}

	if req.Valid.Size() {
		f.Size = req.Size
	}

	if req.Valid.Atime() {
		f.Atime = req.Atime
	}

	if req.Valid.Mtime() {
		f.Mtime = req.Mtime
	}

	if req.Valid.Crtime() {
		f.Crtime = req.Crtime
	}

	if req.Valid.Chgtime() {
		//f.Chgtime = req.Chgtime
	}

	if req.Valid.Bkuptime() {
		//f.Bkuptime = req.Bkuptime
	}

	if req.Valid.Flags() {
		f.Flags = req.Flags
	}
	//_, err := f.namespace.cache.Update(f.namespace.ID, f.FullPath(), f)
	return nil
}

// RemotePath will return the full path on bucket
func (f *File) RemotePath() string {
	return path.Join(f.Parent.RemotePath(), f.Path)
}

// FullPath will return the full path
func (f *File) FullPath() string {
	if !f.IsDirectory() {
		return filepath.Join(f.Parent.FullPath(), f.Path)
	}
	fullPath := ""

	p := f
	for {
		if p == nil {
			break
		}

		fullPath = filepath.Join(p.Path, fullPath)

		p = p.Parent
	}

	return fullPath
}

// Getattr returns the file attributes
func (f *File) Getattr(ctx context.Context) (Attr, error) {
	return f.Attr, nil
}
