package vfs

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

// partly adapted from https://github.com/mark3labs/mcp-filesystem-server/tree/main/filesystemserver/handler
// https://github.com/mark3labs/mcp-filesystem-server/tree/main

// Local fs is a workspace
type LocalFS struct {
	root string

	allowedDirs []string
}

func NewLocalFS(root string) Workspace {
	root = filepath.Clean(root)
	return &LocalFS{
		root:        root,
		allowedDirs: []string{root},
	}
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

// func (s *LocalFS) validatePath(path string) (string, error) {
// 	path = filepath.Clean(path)
// 	rel := strings.TrimPrefix(path, s.root)
// 	abs, err := filepath.Abs(filepath.Join(s.root, rel))
// 	if err != nil {
// 		return "", fmt.Errorf("invalid path %q: %w", path, err)
// 	}
// 	return abs, nil
// }

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

func (s *LocalFS) OpenFile(path string, flag int, perm fs.FileMode) (*os.File, error) {
	validPath, err := s.validatePath(path)
	if err != nil {
		return nil, err
	}
	return os.OpenFile(validPath, flag, perm)
}

func (s *LocalFS) ReadDir(path string) ([]fs.DirEntry, error) {
	validPath, err := s.validatePath(path)
	if err != nil {
		return nil, err
	}
	return os.ReadDir(validPath)
}

func (s *LocalFS) Lstat(path string) (fs.FileInfo, error) {
	validPath, err := s.validatePath(path)
	if err != nil {
		return nil, err
	}
	return os.Lstat(validPath)
}

func (s *LocalFS) Stat(path string) (fs.FileInfo, error) {
	validPath, err := s.validatePath(path)
	if err != nil {
		return nil, err
	}
	return os.Stat(validPath)
}
