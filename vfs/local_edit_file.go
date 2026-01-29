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
	o *EditOptions,
) (int, error) {
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
	if o.UseRegex {
		re, err := regexp.Compile(o.Find)
		if err != nil {
			return -1, err
		}

		if o.AllOccurrences {
			modifiedContent = re.ReplaceAllString(originalContent, o.Replace)
			replacementCount = len(re.FindAllString(originalContent, -1))
		} else {
			matched := re.FindStringIndex(originalContent)
			if matched != nil {
				replacementCount = 1
				modifiedContent = originalContent[:matched[0]] + o.Replace + originalContent[matched[1]:]
			} else {
				modifiedContent = originalContent
				replacementCount = 0
			}
		}
	} else {
		if o.AllOccurrences {
			replacementCount = strings.Count(originalContent, o.Find)
			modifiedContent = strings.ReplaceAll(originalContent, o.Find, o.Replace)
		} else {
			if index := strings.Index(originalContent, o.Find); index != -1 {
				replacementCount = 1
				modifiedContent = originalContent[:index] + o.Replace + originalContent[index+len(o.Find):]
			} else {
				modifiedContent = originalContent
				replacementCount = 0
			}
		}
	}

	if err := os.WriteFile(validPath, []byte(modifiedContent), 0644); err != nil {
		return -1, err
	}

	return replacementCount, nil
}
