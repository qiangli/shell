package vfs

import (
	"fmt"
	"os"
	"path/filepath"
)

func (s *LocalFS) MoveFile(source, destination string) error {
	// Handle empty or relative paths for source
	if source == "." || source == "./" {
		// Get current working directory
		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("Error resolving current directory: %v", err)
		}
		source = cwd
	}

	// Handle empty or relative paths for destination
	if destination == "." || destination == "./" {
		// Get current working directory
		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("Error resolving current directory: %v", err)
		}
		destination = cwd
	}

	validSource, err := s.validatePath(source)
	if err != nil {
		return fmt.Errorf("Error with source path: %v", err)
	}

	// Check if source exists
	if _, err := os.Stat(validSource); os.IsNotExist(err) {
		return fmt.Errorf("Error: Source does not exist: %s", source)
	}

	// For destination path, validate the parent directory first and create it if needed
	destDir := filepath.Dir(destination)
	validDestDir, err := s.validatePath(destDir)
	if err != nil {
		return fmt.Errorf("Error with destination directory path: %v", err)
	}

	// Create parent directory for destination if it doesn't exist
	if err := os.MkdirAll(validDestDir, 0755); err != nil {
		return fmt.Errorf("Error creating destination directory: %v", err)
	}

	// Now validate the full destination path
	validDest, err := s.validatePath(destination)
	if err != nil {
		return fmt.Errorf("Error with destination path: %v", err)
	}

	if err := os.Rename(validSource, validDest); err != nil {
		return fmt.Errorf("Error moving file: %v", err)
	}

	// resourceURI := pathToResourceURI(validDest)
	return nil
}
