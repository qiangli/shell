package vfs

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
	"strings"
)

const (
	maxLineLen   = 80 // Define this based on your Python configuration
	lineNumWidth = 4  // Define this based on your Python configuration
)

type ReadOptions struct {
	// Number the output lines, starting at 1. (cat -n style).
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

	offset := max(o.Offset, 0)
	limit := o.Limit
	if limit <= 0 {
		limit = math.MaxInt
	}

	lines, err := readFile(validPath, o.Number, offset, limit)
	if err != nil {
		return nil, err
	}
	return []byte(lines), nil
}

func readFile(path string, number bool, offset int, limit int) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("error opening file '%s': %v", path, err)
	}
	defer file.Close()

	return ReadLines(file, number, offset, limit)
}

// Read and format content with line numbers
func ReadLines(reader io.Reader, number bool, offset int, limit int) (string, error) {
	scanner := bufio.NewScanner(reader)
	var content []string
	for scanner.Scan() {
		content = append(content, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	if len(content) == 0 {
		return "Empty content", nil
	}

	startIdx := offset
	endIdx := startIdx + limit
	if startIdx >= len(content) {
		return "", fmt.Errorf("error: line offset %d exceeds file length (%d lines)", offset, len(content))
	}

	if endIdx > len(content) {
		endIdx = len(content)
	}

	lines := content[startIdx:endIdx]
	if !number {
		return strings.Join(lines, "\n"), nil
	}
	return FormatLinesWithLineNumbers(lines, startIdx+1), nil
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
