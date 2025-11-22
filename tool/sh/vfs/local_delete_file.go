package vfs

import (
	"fmt"
	"os"
)

func (s *LocalFS) DeleteFile(
	path string,
	recursive bool,
) error {
	validPath, err := s.validatePath(path)
	if err != nil {
		return err
	}

	// Check if path exists
	info, err := os.Stat(validPath)
	if os.IsNotExist(err) {
		return fmt.Errorf("Error: Path does not exist: %s", path)
	} else if err != nil {
		return err
	}

	// Check if it's a directory and handle accordingly
	if info.IsDir() {
		if !recursive {
			return fmt.Errorf("Error: %s is a directory. Use recursive=true to delete directories.", path)
		}

		// It's a directory and recursive is true, so remove it
		if err := os.RemoveAll(validPath); err != nil {
			return err
		}

		return nil
	}

	// It's a file, delete it
	if err := os.Remove(validPath); err != nil {
		return err
	}

	return nil
}
