package vfs

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

func (s *LocalFS) EditFile(
	path string,
	find string,
	replace string,
) (string, error) {
	// Extract optional arguments with defaults
	allOccurrences := true // Default value
	// if val, ok := options["all_occurrences"]; ok {
	// 	allOccurrences = val
	// }

	useRegex := false // Default value
	// if val, ok := options["regex"]; ok {
	// 	useRegex = val
	// }

	// Handle empty or relative paths like "." or "./" by converting to absolute path
	// if path == "." || path == "./" {
	// 	// Get current working directory
	// 	cwd, err := os.Getwd()
	// 	if err != nil {
	// 		return "", err
	// 	}
	// 	path = cwd
	// }

	// Validate path is within allowed directories
	validPath, err := s.validatePath(path)
	if err != nil {
		return "", err
	}

	// Check if it's a directory
	if info, err := os.Stat(validPath); err == nil && info.IsDir() {
		return "", err
	}

	// Check if file exists
	if _, err := os.Stat(validPath); os.IsNotExist(err) {
		return "", err
	}

	// Read file content
	content, err := os.ReadFile(validPath)
	if err != nil {
		return "", err
	}

	originalContent := string(content)
	modifiedContent := ""
	replacementCount := 0

	// Perform the replacement
	if useRegex {
		re, err := regexp.Compile(find)
		if err != nil {
			return "", err
		}

		if allOccurrences {
			modifiedContent = re.ReplaceAllString(originalContent, replace)
			replacementCount = len(re.FindAllString(originalContent, -1))
		} else {
			matched := re.FindStringIndex(originalContent)
			if matched != nil {
				replacementCount = 1
				modifiedContent = originalContent[:matched[0]] + replace + originalContent[matched[1]:]
			} else {
				modifiedContent = originalContent
				replacementCount = 0
			}
		}
	} else {
		if allOccurrences {
			replacementCount = strings.Count(originalContent, find)
			modifiedContent = strings.ReplaceAll(originalContent, find, replace)
		} else {
			if index := strings.Index(originalContent, find); index != -1 {
				replacementCount = 1
				modifiedContent = originalContent[:index] + replace + originalContent[index+len(find):]
			} else {
				modifiedContent = originalContent
				replacementCount = 0
			}
		}
	}

	// Write modified content back to file
	if err := os.WriteFile(validPath, []byte(modifiedContent), 0644); err != nil {
		return "", err
	}

	// Get file info for the response
	info, err := os.Stat(validPath)
	if err != nil {
		// File was written but we couldn't get info
		return fmt.Sprintf("File modified successfully. Made %d replacement(s).", replacementCount), nil
	}

	return fmt.Sprintf("File modified successfully. Made %d replacement(s) in %s (file size: %d bytes)",
		replacementCount, path, info.Size()), nil
}
