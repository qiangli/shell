package sh

import (
	"os"
	"strings"
)

// ClearAllEnv clears all environment variables execep for the keeps
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
