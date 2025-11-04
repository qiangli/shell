package sh

import (
	"context"
	"fmt"
	"os"
	"strings"

	"golang.org/x/exp/slices"
	// cu "github.com/qiangli/shell/tool/coreutils"
	// "github.com/qiangli/shell/tool/coreutils/cat"
	// "github.com/qiangli/shell/tool/sh/vfs"
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

var CoreUtilsCommand = []string{
	"base64", "cat", "chmod", "cp", "find", "gzip", "ls", "mkdir",
	"mktemp", "mv", "rm", "shasum", "tar", "touch", "xargs",
}

// return false without error if not a coreutils
func RunCoreUtils(ctx context.Context, ioe *IOE, args []string) (bool, error) {
	if !slices.Contains(CoreUtilsCommand, args[0]) {
		return false, nil
	}

	switch args[0] {
	case "cat":
		// cmd := cat.New()
		// cmd.SetIO(ioe.Stdin, ioe.Stdout, ioe.Stderr)
		// err := cmd.RunContext(ctx, args[1:]...)
		// if err != nil {
		// 	return true, err
		// }
	default:
		return true, fmt.Errorf("not supported %s", args[0])
	}
	return false, nil
}
