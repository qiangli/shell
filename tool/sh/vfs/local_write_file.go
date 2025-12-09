package vfs

import (
	"fmt"
	"os"
	"path/filepath"
)

func (s *LocalFS) WriteFile(path string, content []byte) error {
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

	// Check if it's a directory
	if info, err := os.Stat(validPath); err == nil && info.IsDir() {
		return fmt.Errorf("Error: Cannot write to a directory")
	}

	// Create parent directories if they don't exist
	parentDir := filepath.Dir(validPath)
	if err := os.MkdirAll(parentDir, 0755); err != nil {
		return fmt.Errorf("Error creating parent directories: %v", err)
	}

	if err := os.WriteFile(validPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("Error writing file: %v", err)
	}

	return nil
}
