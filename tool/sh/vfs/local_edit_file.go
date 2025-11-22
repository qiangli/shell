package vfs

import (
	"os"
	"regexp"
	"strings"
)

type EditOptions struct {
	Find           string
	Replace        string
	AllOccurrences bool
	UseRegex       bool
}

// Update file by finding and replacing text using string matching or regex
// Parameters:
// path (required): Path to the file to modify,
// find (required): Text to search for,
// replace (required): Text to replace with,
// all_occurrences (optional): Replace all occurrences (default: true),
// regex (optional): Treat find pattern as regex (default: false)
// Return replacement count.
func (s *LocalFS) EditFile(
	path string,
	option *EditOptions,
) (int, error) {
	// allOccurrences := true // Default value
	// useRegex := false      // Default value

	// Validate path is within allowed directories
	validPath, err := s.validatePath(path)
	if err != nil {
		return -1, err
	}

	// Check if it's a directory
	if info, err := os.Stat(validPath); err == nil && info.IsDir() {
		return -1, err
	}

	// Check if file exists
	if _, err := os.Stat(validPath); os.IsNotExist(err) {
		return -1, err
	}

	// Read file content
	content, err := os.ReadFile(validPath)
	if err != nil {
		return -1, err
	}

	originalContent := string(content)
	modifiedContent := ""
	replacementCount := 0

	// Perform the replacement
	if option.UseRegex {
		re, err := regexp.Compile(option.Find)
		if err != nil {
			return -1, err
		}

		if option.AllOccurrences {
			modifiedContent = re.ReplaceAllString(originalContent, option.Replace)
			replacementCount = len(re.FindAllString(originalContent, -1))
		} else {
			matched := re.FindStringIndex(originalContent)
			if matched != nil {
				replacementCount = 1
				modifiedContent = originalContent[:matched[0]] + option.Replace + originalContent[matched[1]:]
			} else {
				modifiedContent = originalContent
				replacementCount = 0
			}
		}
	} else {
		if option.AllOccurrences {
			replacementCount = strings.Count(originalContent, option.Find)
			modifiedContent = strings.ReplaceAll(originalContent, option.Find, option.Replace)
		} else {
			if index := strings.Index(originalContent, option.Find); index != -1 {
				replacementCount = 1
				modifiedContent = originalContent[:index] + option.Replace + originalContent[index+len(option.Find):]
			} else {
				modifiedContent = originalContent
				replacementCount = 0
			}
		}
	}

	// Write modified content back to file
	if err := os.WriteFile(validPath, []byte(modifiedContent), 0644); err != nil {
		return -1, err
	}

	// // Get file info for the response
	// info, err := os.Stat(validPath)
	// if err != nil {
	// 	// File was written but we couldn't get info
	// 	return fmt.Sprintf("File modified successfully. Made %d replacement(s).", replacementCount), nil
	// }

	// return fmt.Sprintf("File modified successfully. Made %d replacement(s) in %s (file size: %d bytes)",
	// 	replacementCount, path, info.Size()), nil

	return replacementCount, nil
}
