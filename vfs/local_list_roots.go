package vfs

// import (
// 	// "fmt"
// 	// "path/filepath"
// 	// "strings"
// )

// List allowed top level directories
func (s *LocalFS) ListRoots() ([]string, error) {
	// // Remove the trailing separator for display purposes
	// displayDirs := make([]string, len(s.allowedDirs))
	// for i, dir := range s.allowedDirs {
	// 	displayDirs[i] = strings.TrimSuffix(dir, string(filepath.Separator))
	// }

	// var result strings.Builder
	// result.WriteString("Allowed directories:\n\n")

	// for _, dir := range displayDirs {
	// 	resourceURI := pathToResourceURI(dir)
	// 	result.WriteString(fmt.Sprintf("%s (%s)\n", dir, resourceURI))
	// }

	return s.allowedDirs, nil
}
