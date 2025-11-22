package vfs

import (
	"bytes"
	"fmt"
	"mime"
	"path/filepath"
	"runtime"
	"slices"
	"strconv"
	"strings"

	"github.com/gabriel-vasile/mimetype"
	pt "github.com/qiangli/filesearch"
)

func init() {
	if cpu := runtime.NumCPU(); cpu == 1 {
		runtime.GOMAXPROCS(2)
	} else {
		runtime.GOMAXPROCS(cpu)
	}
}

// Usage:
//   pt [OPTIONS] PATTERN [PATH]

// Application Options:
//       --version             Show version

// Output Options:
//       --color               Print color codes in results (default: true)
//       --nocolor             Don't print color codes in results (default:
//                             false)
//       --color-line-number=  Color codes for line numbers (default: 1;33)
//       --color-path=         Color codes for path names (default: 1;32)
//       --color-match=        Color codes for result matches (default: 30;43)
//       --group               Print file name at header (default: true)
//       --nogroup             Don't print file name at header (default:
//                             false)
//   -0, --null                Separate filenames with null (for 'xargs -0')
//                             (default: false)
//       --column              Print column (default: false)
//       --numbers             Print Line number. (default: true)
//   -N, --nonumbers           Omit Line number. (default: false)
//   -A, --after=              Print lines after match
//   -B, --before=             Print lines before match
//   -C, --context=            Print lines before and after match
//   -l, --files-with-matches  Only print filenames that contain matches
//   -c, --count               Only print the number of matching lines for
//                             each input file.
//   -o, --output-encode=      Specify output encoding (none, jis, sjis, euc)

// Search Options:
//   -e                        Parse PATTERN as a regular expression
//                             (default: false). Accepted syntax is the same
//                             as https://github.com/google/re2/wiki/Syntax
//                             except from \C
//   -i, --ignore-case         Match case insensitively
//   -S, --smart-case          Match case insensitively unless PATTERN
//                             contains uppercase characters
//   -w, --word-regexp         Only match whole words
//       --ignore=             Ignore files/directories matching pattern
//       --vcs-ignore=         VCS ignore files (default: .gitignore)
//       --global-gitignore    Use git's global gitignore file for ignore
//                             patterns
//       --home-ptignore       Use $Home/.ptignore file for ignore patterns
//   -U, --skip-vcs-ignores    Don't use VCS ignore file for ignore patterns
//   -g=                       Print filenames matching PATTERN
//   -G, --file-search-regexp= PATTERN Limit search to filenames matching
//                             PATTERN
//       --depth=              Search up to NUM directories deep (default: 25)
//   -f, --follow              Follow symlinks
//       --hidden              Search hidden files and directories

// Help Options:
//   -h, --help                Show this help message

func Search(path string, o *SearchOptions) (string, error) {
	var args = []string{}

	// output options
	args = append(args, "--output-encode=none")

	// search options
	if o.Regexp {
		args = append(args, "-e")
	}
	if o.IgnoreCase {
		args = append(args, "--ignore-case")
	}
	if o.WordRegexp {
		args = append(args, "--word-regexp")
	}
	for _, exclude := range o.Exclude {
		args = append(args, "--ignore="+exclude)
	}
	if o.FileSearchRegexp != "" {
		args = append(args, "--file-search-regexp="+o.FileSearchRegexp)
	}
	if o.Depth > 0 {
		args = append(args, "--depth="+strconv.Itoa(o.Depth))
	}
	if o.Follow {
		args = append(args, "--follow")
	}
	if o.Hidden {
		args = append(args, "--hidden")
	}

	// args
	if o.Pattern == "" {
		return "", fmt.Errorf("file serach pattern is required")
	}
	args = append(args, o.Pattern)
	if path != "" {
		args = append(args, path)
	}

	var stdoutBuf, stderrBuf bytes.Buffer
	pt := pt.PlatinumSearcher{Out: &stdoutBuf, Err: &stderrBuf}
	exitCode := pt.Run(args)

	if exitCode != 0 {
		return stderrBuf.String(), fmt.Errorf("search failed with exit code %d", exitCode)
	}

	return stdoutBuf.String(), nil
}

// detectMimeType tries to determine the MIME type of a file
func detectMimeType(path string) string {
	// Use mimetype library for more accurate detection
	mtype, err := mimetype.DetectFile(path)
	if err != nil {
		// Fallback to extension-based detection if file can't be read
		ext := filepath.Ext(path)
		if ext != "" {
			mimeType := mime.TypeByExtension(ext)
			if mimeType != "" {
				return mimeType
			}
		}
		return "application/octet-stream" // Default
	}

	return mtype.String()
}

// isTextFile determines if a file is likely a text file based on MIME type
func isTextFile(mimeType string) bool {
	// Check for common text MIME types
	if strings.HasPrefix(mimeType, "text/") {
		return true
	}

	// Common application types that are text-based
	textApplicationTypes := []string{
		"application/json",
		"application/xml",
		"application/javascript",
		"application/x-javascript",
		"application/typescript",
		"application/x-typescript",
		"application/x-yaml",
		"application/yaml",
		"application/toml",
		"application/x-sh",
		"application/x-shellscript",
	}

	if slices.Contains(textApplicationTypes, mimeType) {
		return true
	}

	// Check for +format types
	if strings.Contains(mimeType, "+xml") ||
		strings.Contains(mimeType, "+json") ||
		strings.Contains(mimeType, "+yaml") {
		return true
	}

	// Common code file types that might be misidentified
	if strings.HasPrefix(mimeType, "text/x-") {
		return true
	}

	if strings.HasPrefix(mimeType, "application/x-") &&
		(strings.Contains(mimeType, "script") ||
			strings.Contains(mimeType, "source") ||
			strings.Contains(mimeType, "code")) {
		return true
	}

	return false
}

// isImageFile determines if a file is an image based on MIME type
func isImageFile(mimeType string) bool {
	return strings.HasPrefix(mimeType, "image/") ||
		(mimeType == "application/xml" && strings.HasSuffix(strings.ToLower(mimeType), ".svg"))
}
