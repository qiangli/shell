package vfs

import (
	"fmt"
	"io/fs"
	"os"
	"time"
)

// FileSystem is a virtual file system that provides a set of operations
// to interact with the file system in a controlled manner.
type FileSystem interface {
	FileStore

	ListDirectory(string) ([]string, error)
	CreateDirectory(string) error
	RenameFile(string, string) error
	FileInfo(string) (*FileInfo, error)
	// DeleteFile(string) error
	CopyFile(string, string) error
	EditFile(string, string, string) (string, error)
	Tree(string, int, bool) (string, error)
	SearchFiles(pattern string, path string, options *SearchOptions) (string, error)
}

type FileStore interface {
	ReadFile(string) ([]byte, error)
	WriteFile(string, []byte) error

	// aka:
	// absolute path for file, endpoint for rest, and url for web
	Locator(string) (string, error)
}

type FileStat interface {
	Lstat(name string) (fs.FileInfo, error)
	Stat(name string) (fs.FileInfo, error)
}

type Workspace interface {
	FileSystem

	OpenFile(name string, flag int, perm fs.FileMode) (*os.File, error)
	ReadDir(name string) ([]fs.DirEntry, error)
}

type SearchOptions struct {
	// Parse PATTERN as a regular expression
	// Accepted syntax is the same
	// as https://github.com/google/re2/wiki/Syntax
	Regexp bool
	// Match case insensitively
	IgnoreCase bool
	// Only match whole words
	WordRegexp bool
	// Ignore files/directories matching pattern
	Exclude []string
	// Limit search to filenames matching PATTERN
	FileSearchRegexp string
	// Search up to 'Depth' directories deep (default: 25)
	Depth int
	// Follow symlinks
	Follow bool
	// Search hidden files and directories
	Hidden bool
}

type FileInfo struct {
	Filename string `json:"filename"`

	IsDirectory bool `json:"isDirectory"`
	IsFile      bool `json:"isFile"`
	IsLink      bool `json:"isLink"`

	Permissions string `json:"permissions"`

	Length   int64     `json:"size"`
	Created  time.Time `json:"created"`
	Modified time.Time `json:"modified"`
	Accessed time.Time `json:"accessed"`

	// original info for local fs
	Info fs.FileInfo `json:"-"`
}

// base name of the file
func (r *FileInfo) Name() string {
	return r.Filename
}

// length in bytes for regular files; system-dependent for others
func (r *FileInfo) Size() int64 {
	return r.Length
}

// file mode bits
func (r *FileInfo) Mode() fs.FileMode {
	if r.Info != nil {
		return r.Info.Mode()
	}

	// construct if not set
	var m fs.FileMode
	if r.IsLink {
		m |= fs.ModeSymlink
	}
	if r.IsDirectory {
		m |= fs.ModeDir
	}
	if r.IsFile {
		// specify other flags as needed for regular files.
		m |= 700
	}
	return m
}

// modification time
func (r *FileInfo) ModTime() time.Time {
	return r.Modified
}

// abbreviation for Mode().IsDir()
func (r *FileInfo) IsDir() bool {
	return r.IsDirectory
}

// underlying data source (can return nil)
func (r *FileInfo) Sys() any {
	return nil
}

func (f *FileInfo) String() string {
	return fmt.Sprintf(
		"IsDirectory: %t, IsFile: %t, IsSymlink: %t, Permissions: %s, Size: %d, Created: %s, Modified: %s, Accessed: %s",
		f.IsDirectory,
		f.IsFile,
		f.IsLink,
		f.Permissions,
		f.Length,
		f.Created.Format(time.RFC3339),
		f.Modified.Format(time.RFC3339),
		f.Accessed.Format(time.RFC3339),
	)
}
