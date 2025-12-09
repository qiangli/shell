package vfs

import (
	"fmt"
	"os"
)

func (s *LocalFS) CreateDirectory(path string) error {
	// Handle empty or relative paths like "." or "./" by converting to absolute path
	if path == "." || path == "./" {
		// Get current working directory
		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("Error resolving current directory: %v", err)
		}
		path = cwd
	}

	validPath, err := s.validatePath(path)
	if err != nil {
		return err
	}

	// Check if path already exists
	if info, err := os.Stat(validPath); err == nil {
		if info.IsDir() {
			return fmt.Errorf("Directory already exists: %s", path)
		}
		return fmt.Errorf("Error: Path exists but is not a directory: %s", path)
	}

	if err := os.MkdirAll(validPath, 0755); err != nil {
		return fmt.Errorf("Error creating directory: %v", err)
	}

	return nil
}
