package sh

import (
	"context"
	"os"

	"github.com/qiangli/shell/tool/coreutils/core/backoff"
	"github.com/qiangli/shell/tool/coreutils/core/basename"
	"github.com/qiangli/shell/tool/coreutils/core/cat"
	"github.com/qiangli/shell/tool/coreutils/core/date"
	"github.com/qiangli/shell/tool/coreutils/core/dirname"
	"github.com/qiangli/shell/tool/coreutils/core/head"
	"github.com/qiangli/shell/tool/coreutils/core/tail"
	"github.com/qiangli/shell/tool/coreutils/core/wget"
	"github.com/qiangli/shell/tool/coreutils/exp/tac"

	"github.com/u-root/u-root/pkg/core"
	"github.com/u-root/u-root/pkg/core/base64"
	"github.com/u-root/u-root/pkg/core/chmod"
	"github.com/u-root/u-root/pkg/core/cp"
	"github.com/u-root/u-root/pkg/core/find"
	"github.com/u-root/u-root/pkg/core/gzip"
	"github.com/u-root/u-root/pkg/core/ls"
	"github.com/u-root/u-root/pkg/core/mkdir"
	"github.com/u-root/u-root/pkg/core/mktemp"
	"github.com/u-root/u-root/pkg/core/mv"
	"github.com/u-root/u-root/pkg/core/rm"
	"github.com/u-root/u-root/pkg/core/shasum"
	"github.com/u-root/u-root/pkg/core/tar"
	"github.com/u-root/u-root/pkg/core/touch"
	"github.com/u-root/u-root/pkg/core/xargs"
	"golang.org/x/exp/slices"
)

// tool/coreutils/core/cmp/
//
// tool/coreutils/core/grep/
// tool/coreutils/core/md5sum/
// tool/coreutils/core/seq/
// tool/coreutils/core/sleep/
// tool/coreutils/core/sort/
// tool/coreutils/core/tee/
// tool/coreutils/core/time/
// tool/coreutils/core/timeout/
// tool/coreutils/core/truncate/
// tool/coreutils/core/uniq/
// tool/coreutils/core/wc/
//

// internal commands
var CoreUtilsCommands = []string{
	"base64", "basename", "cat", "chmod", "cp", "date", "dirname", "find", "gzip", "head", "ls", "mkdir",
	"mktemp", "mv", "rm", "shasum", "tac", "tail", "tar", "touch", "wget", "xargs",
}

// bash commands
var BuiltinCommands = []string{
	"true", "false", "exit", "set", "shift", "unset",
	"echo", "printf", "break", "continue", "pwd", "cd",
	"wait", "builtin", "trap", "type", "source", ".", "command",
	"dirs", "pushd", "popd", "umask", "alias", "unalias",
	"fg", "bg", "getopts", "eval", "test", "[", "exec",
	"return", "read", "mapfile", "readarray", "shopt",
}

func IsCoreUtils(s string) bool {
	return !slices.Contains(CoreUtilsCommands, s)
}

func RunBackoff(ctx context.Context, vs *VirtualSystem, args []string) (bool, error) {
	cb := func(args []string) error {
		if IsCoreUtils(args[0]) {
			_, err := RunCoreUtils(ctx, vs, args)
			return err
		}
		// TODO agent/tool support
		// external commands
		cmd := vs.System.Command(args[0], args[1:]...)
		cmd.Stdin, cmd.Stdout, cmd.Stderr = vs.IOE.Stdin, vs.IOE.Stdout, vs.IOE.Stderr
		cmd.Args = args
		cmd.Env = vs.System.Env()
		cmd.Dir, _ = vs.System.Getwd()
		err := cmd.Run()
		return err
	}
	cmd := backoff.New(cb)
	cmd.SetIO(vs.IOE.Stdin, vs.IOE.Stdout, vs.IOE.Stderr)
	err := cmd.RunContext(ctx, args[1:]...)
	return true, err
}

func RunCoreUtils(ctx context.Context, vs *VirtualSystem, args []string) (bool, error) {
	runCmd := func(cmd core.Command) (bool, error) {
		cmd.SetIO(vs.IOE.Stdin, vs.IOE.Stdout, vs.IOE.Stderr)
		cmd.SetWorkingDir(vs.Root)
		err := cmd.RunContext(ctx, args[1:]...)
		return true, err
	}

	open := func(s string) (*os.File, error) {
		return vs.Workspace.OpenFile(s, os.O_RDWR, 0o755)
	}

	switch args[0] {
	case "base64":
		return runCmd(base64.New())
	case "basename":
		return runCmd(basename.New())
	case "cat":
		return runCmd(cat.New(open))
	case "chmod":
		return runCmd(chmod.New())
	case "cp":
		return runCmd(cp.New())
	case "date":
		return runCmd(date.New())
	case "dirname":
		return runCmd(dirname.New())
	case "find":
		return runCmd(find.New())
	case "gzip":
		return runCmd(gzip.New())
	case "head":
		return runCmd(head.New(open))
	case "ls":
		return runCmd(ls.New())
	case "mkdir":
		return runCmd(mkdir.New())
	case "mktemp":
		return runCmd(mktemp.New())
	case "mv":
		return runCmd(mv.New())
	case "rm":
		return runCmd(rm.New())
	case "shasum":
		return runCmd(shasum.New())
	case "tac":
		return runCmd(tac.New(open))
	case "tail":
		return runCmd(tail.New(open))
	case "tar":
		return runCmd(tar.New())
	case "touch":
		return runCmd(touch.New())
	case "wget":
		return runCmd(wget.New())
	case "xargs":
		return runCmd(xargs.New())
	default:
		return false, nil
	}
}
