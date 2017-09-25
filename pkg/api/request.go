package api

import (
	"fmt"
	"os"
	"time"
)

// The SetattrValid are bit flags describing which fields in the SetattrRequest
// are included in the change.
type SetattrValid uint32

const (
	SetattrMode   SetattrValid = 1 << 0
	SetattrUid    SetattrValid = 1 << 1
	SetattrGid    SetattrValid = 1 << 2
	SetattrSize   SetattrValid = 1 << 3
	SetattrAtime  SetattrValid = 1 << 4
	SetattrMtime  SetattrValid = 1 << 5
	SetattrHandle SetattrValid = 1 << 6

	// Linux only(?)
	SetattrAtimeNow  SetattrValid = 1 << 7
	SetattrMtimeNow  SetattrValid = 1 << 8
	SetattrLockOwner SetattrValid = 1 << 9 // http://www.mail-archive.com/git-commits-head@vger.kernel.org/msg27852.html

	// OS X only
	SetattrCrtime   SetattrValid = 1 << 28
	SetattrChgtime  SetattrValid = 1 << 29
	SetattrBkuptime SetattrValid = 1 << 30
	SetattrFlags    SetattrValid = 1 << 31
)

func (fl SetattrValid) Mode() bool      { return fl&SetattrMode != 0 }
func (fl SetattrValid) Uid() bool       { return fl&SetattrUid != 0 }
func (fl SetattrValid) Gid() bool       { return fl&SetattrGid != 0 }
func (fl SetattrValid) Size() bool      { return fl&SetattrSize != 0 }
func (fl SetattrValid) Atime() bool     { return fl&SetattrAtime != 0 }
func (fl SetattrValid) Mtime() bool     { return fl&SetattrMtime != 0 }
func (fl SetattrValid) Handle() bool    { return fl&SetattrHandle != 0 }
func (fl SetattrValid) AtimeNow() bool  { return fl&SetattrAtimeNow != 0 }
func (fl SetattrValid) MtimeNow() bool  { return fl&SetattrMtimeNow != 0 }
func (fl SetattrValid) LockOwner() bool { return fl&SetattrLockOwner != 0 }
func (fl SetattrValid) Crtime() bool    { return fl&SetattrCrtime != 0 }
func (fl SetattrValid) Chgtime() bool   { return fl&SetattrChgtime != 0 }
func (fl SetattrValid) Bkuptime() bool  { return fl&SetattrBkuptime != 0 }
func (fl SetattrValid) Flags() bool     { return fl&SetattrFlags != 0 }

type flagName struct {
	bit  uint32
	name string
}

var setattrValidNames = []flagName{
	{uint32(SetattrMode), "SetattrMode"},
	{uint32(SetattrUid), "SetattrUid"},
	{uint32(SetattrGid), "SetattrGid"},
	{uint32(SetattrSize), "SetattrSize"},
	{uint32(SetattrAtime), "SetattrAtime"},
	{uint32(SetattrMtime), "SetattrMtime"},
	{uint32(SetattrHandle), "SetattrHandle"},
	{uint32(SetattrAtimeNow), "SetattrAtimeNow"},
	{uint32(SetattrMtimeNow), "SetattrMtimeNow"},
	{uint32(SetattrLockOwner), "SetattrLockOwner"},
	{uint32(SetattrCrtime), "SetattrCrtime"},
	{uint32(SetattrChgtime), "SetattrChgtime"},
	{uint32(SetattrBkuptime), "SetattrBkuptime"},
	{uint32(SetattrFlags), "SetattrFlags"},
}

func (fl SetattrValid) String() string {
	return flagString(uint32(fl), setattrValidNames)
}

func flagString(f uint32, names []flagName) string {
	var s string

	if f == 0 {
		return "0"
	}

	for _, n := range names {
		if f&n.bit != 0 {
			s += "+" + n.name
			f &^= n.bit
		}
	}
	if f != 0 {
		s += fmt.Sprintf("%+#x", f)
	}
	return s[1:]
}

// A SetattrRequest asks to change one or more attributes associated with a file,
// as indicated by Valid.
type SetattrRequest struct {
	Valid SetattrValid
	Size  uint64
	Atime time.Time
	Mtime time.Time
	Mode  os.FileMode
	Uid   uint32
	Gid   uint32

	// OS X only
	Bkuptime time.Time
	Chgtime  time.Time
	Crtime   time.Time
	Flags    uint32 // see chflags(2)
}

// A GetattrRequest asks for the metadata for the file denoted by r.Node.
type GetattrRequest struct {
	Flags uint32 // bit flags
}

// A LookupRequest asks to look up the given name in the directory named by r.Node.
type LookupRequest struct {
	Name string
}

// OpenFlags are the O_FOO flags passed to open/create/etc calls. For
// example, os.O_WRONLY | os.O_APPEND.
type OpenFlags uint32

// A CreateRequest asks to create and open a file (not a directory).
type CreateRequest struct {
	Name  string
	Flags OpenFlags
	Mode  os.FileMode
	// Umask of the request. Not supported on OS X.
	Umask os.FileMode
}

type MkdirRequest struct {
	Name string
	Mode os.FileMode
	// Umask of the request. Not supported on OS X.
	Umask os.FileMode
}

// The ReleaseFlags are used in the Release exchange.
type ReleaseFlags uint32

// A ReleaseRequest asks to release (close) an open file handle.
type ReleaseRequest struct {
	Dir          bool      // is this Releasedir?
	Flags        OpenFlags // flags from OpenRequest
	ReleaseFlags ReleaseFlags
	LockOwner    uint32
}

// A FlushRequest asks for the current state of an open file to be flushed
// to storage, as when a file descriptor is being closed.  A single opened Handle
// may receive multiple FlushRequests over its lifetime.
type FlushRequest struct {
	Flags     uint32
	LockOwner uint64
}

// A RemoveRequest asks to remove a file or directory from the
// directory r.Node.
type RemoveRequest struct {
	Name string // name of the entry to remove
	Dir  bool   // is this rmdir?
}

// A SymlinkRequest is a request to create a symlink making NewName point to Target.
type SymlinkRequest struct {
	NewName string
	Target  string
}

// A LinkRequest is a request to create a hard link.
type LinkRequest struct {
	OldNode uint64
	NewName string
}

// A RenameRequest is a request to rename a file.
type RenameRequest struct {
	NewDir  uint64
	OldName string
	NewName string
}

type MknodRequest struct {
	Name string
	Mode os.FileMode
	Rdev uint32
	// Umask of the request. Not supported on OS X.
	Umask os.FileMode
}

type FsyncRequest struct {
	// TODO bit 1 is datasync, not well documented upstream
	Flags uint32
	Dir   bool
}
