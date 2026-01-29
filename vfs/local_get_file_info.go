package vfs

import (
	"fmt"
	"os"
	"time"

	"github.com/djherbis/times"
)

func (s *LocalFS) GetFileInfo(path string) (*FileInfo, error) {
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

	info, err := s.getFileStats(validPath)
	if err != nil {
		return nil, fmt.Errorf("Error getting file info: %v", err)
	}
	return info, err

	// Get MIME type for files
	// mimeType := "directory"
	// if info.IsFile {
	// 	mimeType = DetectMimeType(validPath)
	// }
	// return s.getFileStats(validPath)
	// resourceURI := pathToResourceURI(validPath)

	// Determine file type text
	// var fileTypeText string
	// if info.IsDirectory {
	// 	fileTypeText = "Directory"
	// } else {
	// 	fileTypeText = "File"
	// }

	// return &mcp.CallToolResult{
	// 	Content: []mcp.Content{
	// 		mcp.TextContent{
	// 			Type: "text",
	// 			Text: fmt.Sprintf(
	// 				"File information for: %s\n\nSize: %d bytes\nCreated: %s\nModified: %s\nAccessed: %s\nIsDirectory: %v\nIsFile: %v\nPermissions: %s\nMIME Type: %s\nResource URI: %s",
	// 				validPath,
	// 				info.Size,
	// 				info.Created.Format(time.RFC3339),
	// 				info.Modified.Format(time.RFC3339),
	// 				info.Accessed.Format(time.RFC3339),
	// 				info.IsDirectory,
	// 				info.IsFile,
	// 				info.Permissions,
	// 				mimeType,
	// 				resourceURI,
	// 			),
	// 		},
	// 		mcp.EmbeddedResource{
	// 			Type: "resource",
	// 			Resource: mcp.TextResourceContents{
	// 				URI:      resourceURI,
	// 				MIMEType: "text/plain",
	// 				Text: fmt.Sprintf("%s: %s (%s, %d bytes)",
	// 					fileTypeText,
	// 					validPath,
	// 					mimeType,
	// 					info.Size),
	// 			},
	// 		},
	// 	},
	// }, nil
}

// func (fs *LocalFS) getFileStats(path string) (FileInfo, error) {
// 	info, err := os.Stat(path)
// 	if err != nil {
// 		return FileInfo{}, err
// 	}

// 	timespec, err := times.Stat(path)
// 	if err != nil {
// 		return FileInfo{}, fmt.Errorf("failed to get file times: %w", err)
// 	}

// 	createdTime := time.Time{}
// 	if timespec.HasBirthTime() {
// 		createdTime = timespec.BirthTime()
// 	}

// 	return FileInfo{
// 		Size:        info.Size(),
// 		Created:     createdTime,
// 		Modified:    timespec.ModTime(),
// 		Accessed:    timespec.AccessTime(),
// 		IsDirectory: info.IsDir(),
// 		IsFile:      !info.IsDir(),
// 		Permissions: fmt.Sprintf("%o", info.Mode().Perm()),
// 	}, nil
// }

func (s *LocalFS) getFileStats(path string) (*FileInfo, error) {
	info, err := os.Stat(path)
	if err != nil {
		return &FileInfo{}, err
	}
	isLink := (info.Mode() & os.ModeSymlink) != 0

	timespec, err := times.Stat(path)
	if err != nil {
		return &FileInfo{}, fmt.Errorf("failed to get file times: %w", err)
	}

	createdTime := time.Time{}
	if timespec.HasBirthTime() {
		createdTime = timespec.BirthTime()
	}

	mimeType := "directory"
	if info.Mode().IsRegular() {
		mimeType = DetectMimeType(path)
	}

	return &FileInfo{
		IsDirectory: info.IsDir(),
		IsFile:      info.Mode().IsRegular(),
		IsLink:      isLink,
		Permissions: fmt.Sprintf("%o", info.Mode().Perm()),
		Length:      info.Size(),
		Created:     createdTime,
		Modified:    timespec.ModTime(),
		Accessed:    timespec.AccessTime(),
		//
		Mime: mimeType,
		//
		Info: info,
	}, nil
}
