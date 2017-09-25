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

	"github.com/gostor/gofs/pkg/api"
)

// Lookup -
func (dir *File) Lookup(ctx context.Context, name string) error {
	/*
		if !dir.IsDirectory() {
			return nil, nil
		}
		// we are not statting each object here because of performance reasons
		var o interface{} // meta.Object
		if err := dir.mfs.db.View(func(tx *meta.Tx) error {
			b := dir.bucket(tx)
			return b.Get(name, &o)
		}); err == nil {
		} else if meta.IsNoSuchObject(err) {
			return nil, fuse.ENOENT
		} else if err != nil {
			return nil, err
		}

		if file, ok := o.(File); ok {
			file.mfs = dir.mfs
			file.dir = dir
			return &file, nil
		} else if subdir, ok := o.(Dir); ok {
			subdir.mfs = dir.mfs
			subdir.dir = dir
			return &subdir, nil
		}

		return nil, fuse.ENOENT
	*/
	return nil
}

func (dir *File) scan(ctx context.Context) error {
	/*
			if !dir.needsScan() {
				return nil
			}

			tx, err := dir.mfs.db.Begin(true)
			if err != nil {
				return err
			}

			defer tx.Rollback()

			b := dir.bucket(tx)

			objects := map[string]interface{}{}

			// we'll compare the current bucket contents against our cache folder, and update the cache
			if err := b.ForEach(func(k string, o interface{}) error {
				if k[len(k)-1] != '/' {
					objects[k] = &o
				}
				return nil
			}); err != nil {
				return err
			}

			prefix := dir.RemotePath()
			if prefix != "" {
				prefix = prefix + "/"
			}

			// The channel will abort the ListObjectsV2 request.
			doneCh := make(chan struct{})
			defer close(doneCh)

			ch := dir.mfs.api.ListObjectsV2(dir.mfs.config.bucket, prefix, false, doneCh)

		loop:
			for {
				select {
				case <-ctx.Done():
					break loop
				case objInfo, ok := <-ch:
					if !ok {
						break loop
					}
					key := objInfo.Key[len(prefix):]
					baseKey := path.Base(key)

					// object still exists
					objects[baseKey] = nil

					if strings.HasSuffix(key, "/") {
						dir.storeDir(b, tx, baseKey, objInfo)
					} else {
						dir.storeFile(b, tx, baseKey, objInfo)
					}
				}
			}

			// cache housekeeping
			for k, o := range objects {
				if o == nil {
					continue
				}

				// purge from cache
				b.Delete(k)

				if _, ok := o.(Dir); !ok {
					continue
				}

				b.DeleteBucket(k + "/")
			}

			if err := tx.Commit(); err != nil {
				return err
			}

			dir.scanned = true
	*/
	return nil
}

// ReadDirAll will return all files in current dir
func (dir *File) ReadDirAll(ctx context.Context) error {
	/*
		if err := dir.scan(ctx); err != nil {
			return nil, err
		}

		var entries = []fuse.Dirent{}

		// update cache folder with bucket list
		if err := dir.mfs.db.View(func(tx *meta.Tx) error {
			return dir.bucket(tx).ForEach(func(k string, o interface{}) error {
				if file, ok := o.(File); ok {
					file.dir = dir
					entries = append(entries, file.Dirent())
				} else if subdir, ok := o.(Dir); ok {
					subdir.dir = dir
					entries = append(entries, subdir.Dirent())
				} else {
					panic("Could not find type. Try to remove cache.")
				}

				return nil
			})
		}); err != nil {
			return nil, err
		}

		return entries, nil
	*/
	return nil
}

/*
func (dir *File) bucket(tx *meta.Tx) *meta.Bucket {
	// Root folder.
	if dir.dir == nil {
		return tx.Bucket("minio/")
	}

	b := dir.dir.bucket(tx)

	return b.Bucket(dir.Path + "/")
}
*/

// Mkdir will make a new directory below current dir
func (dir *File) Mkdir(ctx context.Context, req *api.MkdirRequest) error {
	/*
		subdir := File{
			Parent:    dir,
			namespace: dir.namespace,

			Path: req.Name,

			Mode: 0770 | os.ModeDir,
			GID:  dir.mfs.config.gid,
			UID:  dir.mfs.config.uid,

			Chgtime: time.Now(),
			Crtime:  time.Now(),
			Mtime: time.Now(),
			Atime: time.Now(),
		}

			tx, err := dir.mfs.db.Begin(true)
			if err != nil {
				return nil, err
			}

			defer tx.Rollback()

			if err := subdir.store(tx); err != nil {
				return nil, err
			}

			// Commit the transaction and check for error.
			if err := tx.Commit(); err != nil {
				return nil, err
			}
	*/

	//return &subdir, nil
	return nil
}

// Remove will delete a file or directory from current directory
func (dir *File) Remove(ctx context.Context, req *api.RemoveRequest) error {
	/*
		if err := dir.mfs.wait(path.Join(dir.FullPath(), req.Name)); err != nil {
			return err
		}

		tx, err := dir.mfs.db.Begin(true)
		if err != nil {
			return err
		}

		defer tx.Rollback()

		b := dir.bucket(tx)

		var o interface{}
		if err := b.Get(req.Name, &o); meta.IsNoSuchObject(err) {
			return fuse.ENOENT
		} else if err != nil {
			return err
		} else if err := b.Delete(req.Name); err != nil {
			return err
		}

		if req.Dir {
			b.DeleteBucket(req.Name + "/")
		}

		if err := dir.mfs.api.RemoveObject(dir.mfs.config.bucket, path.Join(dir.RemotePath(), req.Name)); err != nil {
			return err
		}

		if err := tx.Commit(); err != nil {
			return err
		}
	*/

	return nil
}

/*
// store the dir object in cache
func (dir *File) store(tx *meta.Tx) error {
	// directories will be stored in their parent buckets
	b := dir.dir.bucket(tx)

	subbucketPath := path.Base(dir.Path)
	if _, err := b.CreateBucketIfNotExists(subbucketPath + "/"); err != nil {
		return err
	}

	return b.Put(subbucketPath, dir)
}
*/

// Create will return a new empty file in current dir, if the file is currently locked, it will
// wait for the lock to be freed.
func (dir *File) Create(ctx context.Context, req *api.CreateRequest) error {
	return nil
}

// Rename will rename files
func (dir *File) Rename(ctx context.Context, req *api.RenameRequest) error {
	return nil
}
