package vfs

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
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
	// EditFile
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

// Local fs is a workspace
type LocalFS struct {
}

func NewLocalFS() Workspace {
	return &LocalFS{}
}

func (s *LocalFS) ListDirectory(path string) ([]string, error) {
	validPath, err := s.validatePath(path)
	if err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(validPath)
	if err != nil {
		return nil, err
	}

	var result []string
	for _, entry := range entries {
		prefix := "File"
		if entry.IsDir() {
			prefix = "Direcotory"
		}
		result = append(result, fmt.Sprintf("%s: %s", prefix, entry.Name()))
	}

	return result, nil
}

func (s *LocalFS) CreateDirectory(path string) error {
	validPath, err := s.validatePath(path)
	if err != nil {
		return err
	}

	return os.MkdirAll(validPath, 0755)
}

func (s *LocalFS) RenameFile(source, destination string) error {
	validSource, err := s.validatePath(source)
	if err != nil {
		return err
	}
	validDest, err := s.validatePath(destination)
	if err != nil {
		return err
	}

	return os.Rename(validSource, validDest)
}

func (s *LocalFS) FileInfo(path string) (*FileInfo, error) {
	validPath, err := s.validatePath(path)
	if err != nil {
		return nil, err
	}

	info, err := s.getFileStats(validPath)
	if err != nil {
		return nil, err
	}

	return info, nil
}

func (s *LocalFS) ReadFile(path string) ([]byte, error) {
	validPath, err := s.validatePath(path)
	if err != nil {
		return nil, err
	}

	return os.ReadFile(validPath)
}

func (s *LocalFS) WriteFile(path string, content []byte) error {
	validPath, err := s.validatePath(path)
	if err != nil {
		return err
	}

	return os.WriteFile(validPath, content, 0644)
}

func (s *LocalFS) Locator(path string) (string, error) {
	return s.validatePath(path)
}

func (s *LocalFS) validatePath(path string) (string, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("invalid path %q: %w", path, err)
	}

	return abs, nil
}

func (s *LocalFS) getFileStats(path string) (*FileInfo, error) {
	info, err := os.Stat(path)
	if err != nil {
		return &FileInfo{}, err
	}
	isLink := (info.Mode() & os.ModeSymlink) != 0

	return &FileInfo{
		IsDirectory: info.IsDir(),
		IsFile:      info.Mode().IsRegular(),
		IsLink:      isLink,
		Permissions: fmt.Sprintf("%o", info.Mode().Perm()),
		Length:      info.Size(),
		Created:     info.ModTime(),
		Modified:    info.ModTime(),
		Accessed:    info.ModTime(),
		//
		Info: info,
	}, nil
}

func (s *LocalFS) SearchFiles(pattern string, path string, options *SearchOptions) (string, error) {
	if options == nil {
		options = &SearchOptions{}
	}
	return Search(pattern, path, options)
}

func (s *LocalFS) OpenFile(name string, flag int, perm fs.FileMode) (*os.File, error) {
	return os.OpenFile(name, flag, perm)
}

func (s *LocalFS) ReadDir(name string) ([]fs.DirEntry, error) {
	return os.ReadDir(name)
}

func (s *LocalFS) Lstat(name string) (fs.FileInfo, error) {
	return os.Lstat(name)
}

func (s *LocalFS) Stat(name string) (fs.FileInfo, error) {
	return os.Stat(name)
}
