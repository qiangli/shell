package vfs

import (
	"fmt"
	"os"
	"strings"
)

const (
	MAX_LINE_LENGTH   = 80 // Define this based on your Python configuration
	LINE_NUMBER_WIDTH = 4  // Define this based on your Python configuration
)

// FormatContentWithLineNumbers formats file content with line numbers (cat -n style).
func FormatContentWithLineNumbers(content string, startLine int) string {
	lines := strings.Split(content, "\n")
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}

	var resultLines []string
	for i, line := range lines {
		lineNum := i + startLine

		if len(line) <= MAX_LINE_LENGTH {
			resultLines = append(resultLines, fmt.Sprintf("%*d\t%s", LINE_NUMBER_WIDTH, lineNum, line))
		} else {
			// Split long line into chunks with continuation markers
			numChunks := (len(line) + MAX_LINE_LENGTH - 1) / MAX_LINE_LENGTH
			for chunkIdx := 0; chunkIdx < numChunks; chunkIdx++ {
				start := chunkIdx * MAX_LINE_LENGTH
				end := start + MAX_LINE_LENGTH
				if end > len(line) {
					end = len(line)
				}
				chunk := line[start:end]
				if chunkIdx == 0 {
					resultLines = append(resultLines, fmt.Sprintf("%*d\t%s", LINE_NUMBER_WIDTH, lineNum, chunk))
				} else {
					continuationMarker := fmt.Sprintf("%d.%d", lineNum, chunkIdx)
					resultLines = append(resultLines, fmt.Sprintf("%*s\t%s", LINE_NUMBER_WIDTH, continuationMarker, chunk))
				}
			}
		}
	}

	return strings.Join(resultLines, "\n")
}

func (s *LocalFS) ReadFile(path string) ([]byte, error) {
	validPath, err := s.validatePath(path)
	if err != nil {
		return nil, err
	}

	return os.ReadFile(validPath)
}
