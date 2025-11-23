package vfs

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
)

const (
	maxLineLen   = 80 // Define this based on your Python configuration
	lineNumWidth = 4  // Define this based on your Python configuration
)

type ReadOptions struct {
	//  Number the output lines, starting at 1. (cat -n style).
	Number bool

	// Line offset to start reading from (0-indexed)
	Offset int

	// Maximum number of lines to read. Return the entire file if zero or less
	Limit int
}

// ReadFile read raw bytes or content with line numbers if option is provided.
// path: Absolute or relative file path
// offset: Line offset to start reading from (0-indexed)
// limit: Maximum number of lines to read
// Returns:
// Formatted file content with line numbers.
func (s *LocalFS) ReadFile(path string, o *ReadOptions) ([]byte, error) {
	validPath, err := s.validatePath(path)
	if err != nil {
		return nil, err
	}

	if o == nil {
		return os.ReadFile(validPath)
	}

	offset := o.Offset
	if offset < 0 {
		offset = 0
	}
	limit := o.Limit
	if limit <= 0 {
		limit = math.MaxInt
	}

	lines, err := readLines(path, offset, limit)
	if err != nil {
		return nil, err
	}
	return []byte(lines), nil
}

func readLines(filePath string, offset int, limit int) (string, error) {
	resolvedPath, err := filepath.Abs(filePath)
	if err != nil {
		return "", fmt.Errorf("error resolving file path: %v", err)
	}

	// Attempt to open the file
	file, err := os.Open(resolvedPath)
	if err != nil {
		return "", fmt.Errorf("error opening file '%s': %v", filePath, err)
	}
	defer file.Close()

	// Read file content
	scanner := bufio.NewScanner(file)
	var content []string
	for scanner.Scan() {
		content = append(content, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("error reading file '%s': %v", filePath, err)
	}

	if len(content) == 0 {
		return "Error: Empty file", nil
	}

	startIdx := offset
	endIdx := startIdx + limit
	if startIdx >= len(content) {
		return "", fmt.Errorf("error: line offset %d exceeds file length (%d lines)", offset, len(content))
	}

	if endIdx > len(content) {
		endIdx = len(content)
	}

	selectedLines := content[startIdx:endIdx]
	return FormatLinesWithLineNumbers(selectedLines, startIdx+1), nil
}

// FormatContentWithLineNumbers formats content with line numbers (cat -n style).
func FormatContentWithLineNumbers(content string, startLine int) string {
	lines := strings.Split(content, "\n")
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	return FormatLinesWithLineNumbers(lines, startLine)
}

func FormatLinesWithLineNumbers(lines []string, startLine int) string {
	var resultLines []string
	for i, line := range lines {
		lineNum := i + startLine

		if len(line) <= maxLineLen {
			resultLines = append(resultLines, fmt.Sprintf("%*d\t%s", lineNumWidth, lineNum, line))
		} else {
			// Split long line into chunks with continuation markers
			numChunks := (len(line) + maxLineLen - 1) / maxLineLen
			for chunkIdx := 0; chunkIdx < numChunks; chunkIdx++ {
				start := chunkIdx * maxLineLen
				end := start + maxLineLen
				if end > len(line) {
					end = len(line)
				}
				chunk := line[start:end]
				if chunkIdx == 0 {
					resultLines = append(resultLines, fmt.Sprintf("%*d\t%s", lineNumWidth, lineNum, chunk))
				} else {
					continuationMarker := fmt.Sprintf("%d.%d", lineNum, chunkIdx)
					resultLines = append(resultLines, fmt.Sprintf("%*s\t%s", lineNumWidth, continuationMarker, chunk))
				}
			}
		}
	}

	return strings.Join(resultLines, "\n")
}
