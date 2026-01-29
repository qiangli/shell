package vfs

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// Copy files and directories
// Parameters:
// source (required): Source path of the file or directory,
// destination (required): Destination path
func (s *LocalFS) CopyFile(
	source string,
	destination string,
) error {
	validSource, err := s.validatePath(source)
	if err != nil {
		return err
	}

	// Check if source exists
	srcInfo, err := os.Stat(validSource)
	if os.IsNotExist(err) {
		return fmt.Errorf("Error: Source does not exist: %s", source)
	} else if err != nil {
		return fmt.Errorf("Error accessing source: %v", err)
	}

	validDest, err := s.validatePath(destination)
	if err != nil {
		return fmt.Errorf("Error with destination path: %v", err)
	}

	// Create parent directory for destination if it doesn't exist
	destDir := filepath.Dir(validDest)
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("Error creating destination directory: %v", err)
	}

	// Perform the copy operation based on whether source is a file or directory
	if srcInfo.IsDir() {
		// It's a directory, copy recursively
		if err := copyDir(validSource, validDest); err != nil {
			return fmt.Errorf("Error copying directory: %v", err)
		}
	} else {
		// It's a file, copy directly
		if err := copyFile(validSource, validDest); err != nil {
			return fmt.Errorf("Error copying file: %v", err)
		}
	}

	return nil
}

// copyFile copies a single file from src to dst
func copyFile(src, dst string) error {
	// Open the source file
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	// Create the destination file
	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	// Copy the contents
	if _, err := io.Copy(destFile, sourceFile); err != nil {
		return err
	}

	// Get source file mode
	sourceInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	// Set the same file mode on destination
	return os.Chmod(dst, sourceInfo.Mode())
}

// copyDir recursively copies a directory tree from src to dst
func copyDir(src, dst string) error {
	// Get properties of source dir
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	// Create the destination directory with the same permissions
	if err = os.MkdirAll(dst, srcInfo.Mode()); err != nil {
		return err
	}

	// Read directory entries
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		// Handle symlinks
		if entry.Type()&os.ModeSymlink != 0 {
			// For simplicity, we'll skip symlinks in this implementation
			continue
		}

		// Recursively copy subdirectories or copy files
		if entry.IsDir() {
			if err = copyDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			if err = copyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}

	return nil
}
