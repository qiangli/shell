package sh

import (
	"fmt"
	"io/fs"
	"os"
	"strings"
)

// ClearAllEnv clears all environment variables except for the keeps
func ClearAllEnv(keeps []string) {
	var memo = make(map[string]bool, len(keeps))
	for _, key := range keeps {
		memo[key] = true
	}

	for _, env := range os.Environ() {
		key := strings.Split(env, "=")[0]
		if !memo[key] {
			os.Unsetenv(key)
		}
	}
}

func decodeFileFlag(flag int) string {
	var parts []string
	if flag&os.O_RDONLY != 0 {
		parts = append(parts, "O_RDONLY")
	}
	if flag&os.O_WRONLY != 0 {
		parts = append(parts, "O_WRONLY")
	}
	if flag&os.O_RDWR != 0 {
		parts = append(parts, "O_RDWR")
	}
	if flag&os.O_APPEND != 0 {
		parts = append(parts, "O_APPEND")
	}
	if flag&os.O_CREATE != 0 {
		parts = append(parts, "O_CREATE")
	}
	if flag&os.O_EXCL != 0 {
		parts = append(parts, "O_EXCL")
	}
	if flag&os.O_SYNC != 0 {
		parts = append(parts, "O_SYNC")
	}
	if flag&os.O_TRUNC != 0 {
		parts = append(parts, "O_TRUNC")
	}
	return strings.Join(parts, " | ")
}

func decodeFilePerm(perm fs.FileMode) string {
	return fmt.Sprintf("%#o", perm.Perm())
}
