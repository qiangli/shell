package vfs

import (
	"fmt"
	"os"
	"path/filepath"
)

func (s *LocalFS) ListDirectory(path string) ([]string, error) {
	// Handle empty or relative paths like "." or "./" by converting to absolute path
	if path == "." || path == "./" {
		// Get current working directory
		cwd, err := os.Getwd()
		if err != nil {
			return nil, fmt.Errorf("Error resolving current directory: %v", err)
		}
		path = cwd
	}

	validPath, err := s.validatePath(path)
	if err != nil {
		return nil, err
	}

	// Check if it's a directory
	info, err := os.Stat(validPath)
	if err != nil {
		return nil, err
	}

	if !info.IsDir() {
		return nil, fmt.Errorf("Error: Path is not a directory")
	}

	entries, err := os.ReadDir(validPath)
	if err != nil {
		return nil, fmt.Errorf("Error reading directory: %v", err)
	}

	var result []string
	for _, entry := range entries {
		entryPath := filepath.Join(validPath, entry.Name())
		resourceURI := PathToResourceURI(entryPath)

		if entry.IsDir() {
			result = append(result, fmt.Sprintf("[DIR]  %s (%s)\n", entry.Name(), resourceURI))
		} else {
			info, err := entry.Info()
			if err == nil {
				result = append(result, fmt.Sprintf("[FILE] %s (%s) - %d bytes\n",
					entry.Name(), resourceURI, info.Size()))
			} else {
				result = append(result, fmt.Sprintf("[FILE] %s (%s)\n", entry.Name(), resourceURI))
			}
		}
	}

	return result, nil
}
