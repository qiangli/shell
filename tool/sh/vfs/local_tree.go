package vfs

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// FileNode represents a node in the file tree
type FileNode struct {
	Name     string      `json:"name"`
	Path     string      `json:"path"`
	Type     string      `json:"type"` // "file" or "directory"
	Size     int64       `json:"size,omitempty"`
	Modified time.Time   `json:"modified,omitempty"`
	Children []*FileNode `json:"children,omitempty"`
}

func (s *LocalFS) Tree(
	path string,
	depth int,
	follow bool,
) (string, error) {
	// Validate the path is within allowed directories
	validPath, err := s.validatePath(path)
	if err != nil {
		return "", nil
	}

	// Check if it's a directory
	info, err := os.Stat(validPath)
	if err != nil {
		return "", nil
	}

	if !info.IsDir() {
		return "", err
	}

	// Build the tree structure
	tree, err := s.buildTree(validPath, depth, 0, follow)
	if err != nil {
		return "", nil
	}

	// Convert to JSON
	jsonData, err := json.MarshalIndent(tree, "", "  ")
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Directory tree for %s (max depth: %d):\n\n%s", validPath, depth, string(jsonData)), nil
}

// buildTree builds a tree representation of the filesystem starting at the given path
func (s *LocalFS) buildTree(path string, maxDepth int, currentDepth int, followSymlinks bool) (*FileNode, error) {
	// Validate the path
	validPath, err := s.validatePath(path)
	if err != nil {
		return nil, err
	}

	// Get file info
	info, err := os.Stat(validPath)
	if err != nil {
		return nil, err
	}

	// Create the node
	node := &FileNode{
		Name:     filepath.Base(validPath),
		Path:     validPath,
		Modified: info.ModTime(),
	}

	// Set type and size
	if info.IsDir() {
		node.Type = "directory"

		// If we haven't reached the max depth, process children
		if currentDepth < maxDepth {
			// Read directory entries
			entries, err := os.ReadDir(validPath)
			if err != nil {
				return nil, err
			}

			// Process each entry
			for _, entry := range entries {
				entryPath := filepath.Join(validPath, entry.Name())

				// Handle symlinks
				if entry.Type()&os.ModeSymlink != 0 {
					if !followSymlinks {
						// Skip symlinks if not following them
						continue
					}

					// Resolve symlink
					linkDest, err := filepath.EvalSymlinks(entryPath)
					if err != nil {
						// Skip invalid symlinks
						continue
					}

					// Validate the symlink destination is within allowed directories
					if !s.isPathInAllowedDirs(linkDest) {
						// Skip symlinks pointing outside allowed directories
						continue
					}

					entryPath = linkDest
				}

				// Recursively build child node
				childNode, err := s.buildTree(entryPath, maxDepth, currentDepth+1, followSymlinks)
				if err != nil {
					// Skip entries with errors
					continue
				}

				// Add child to the current node
				node.Children = append(node.Children, childNode)
			}
		}
	} else {
		node.Type = "file"
		node.Size = info.Size()
	}

	return node, nil
}
