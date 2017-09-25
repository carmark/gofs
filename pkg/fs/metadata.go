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
	"os"
	"time"
)

// An Attr is the metadata for a single file or directory.
type Attr struct {
	Inode  uint64      // inode number
	Size   uint64      // size in bytes
	Atime  time.Time   // time of last access
	Mtime  time.Time   // time of last modification
	Ctime  time.Time   // time of last inode change
	Crtime time.Time   // time of creation (OS X only)
	Mode   os.FileMode // file mode
	Nlink  uint32      // number of links (usually 1)
	Uid    uint32      // owner uid
	Gid    uint32      // group gid
	Flags  uint32      // chflags(2) flags (OS X only)
}
