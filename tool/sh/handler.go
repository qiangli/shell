package sh

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"mvdan.cc/sh/v3/expand"
	"mvdan.cc/sh/v3/interp"

	"github.com/qiangli/shell/tool/sh/vfs"
	"github.com/qiangli/shell/tool/sh/vos"
)

// standard IO
type IOE struct {
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

type ExecHandler func(context.Context, []string) (bool, error)

type VirtualSystem struct {
	IOE *IOE

	Workspace vfs.Workspace
	System    vos.System

	ExecHandler ExecHandler

	MaxTimeout int
}

func (s *VirtualSystem) Run(script string) error {
	return runAll(s, script)
}

func NewVirtualSystem(s vos.System, ws vfs.Workspace, ioe *IOE) *VirtualSystem {
	return &VirtualSystem{
		System:    s,
		Workspace: ws,
		IOE:       ioe,
	}
}

func NewLocalSystem(root string, ioe *IOE) *VirtualSystem {
	return NewVirtualSystem(vos.NewLocalSystem(root), vfs.NewLocalFS(root), ioe)
}

func NewDummyExecHandler(ioe *IOE) ExecHandler {
	return func(ctx context.Context, args []string) (bool, error) {
		fmt.Fprintf(ioe.Stderr, "args: %+v\n", args)
		if args[0] == "ai" || strings.HasPrefix(args[0], "@") {
			fmt.Fprintf(ioe.Stdout, "ai args: %+v\n", args)
			return true, nil
		}
		return false, nil
	}
}

func VirtualOpenHandler(ws vfs.Workspace) interp.OpenHandlerFunc {
	return func(ctx context.Context, path string, flag int, perm fs.FileMode) (io.ReadWriteCloser, error) {
		mc := interp.HandlerCtx(ctx)
		//
		if runtime.GOOS == "windows" && path == "/dev/null" {
			path = "NUL"
			// Work around https://go.dev/issue/71752, where Go 1.24 started giving
			// "Invalid handle" errors when opening "NUL" with O_TRUNC.
			// TODO: hopefully remove this in the future once the bug is fixed.
			flag &^= os.O_TRUNC
		} else if path != "" && !filepath.IsAbs(path) {
			path = filepath.Join(mc.Dir, path)
		}
		return ws.OpenFile(path, flag, perm)
	}
}

func VirtualReadDirHandler2(ws vfs.Workspace) interp.ReadDirHandlerFunc2 {
	return func(ctx context.Context, path string) ([]fs.DirEntry, error) {
		return ws.ReadDir(path)
	}
}

func VirtualStatHandler(ws vfs.Workspace) interp.StatHandlerFunc {
	return func(ctx context.Context, path string, followSymlinks bool) (fs.FileInfo, error) {
		if v, ok := ws.(vfs.FileStat); ok {
			if !followSymlinks {
				return v.Lstat(path)
			} else {
				return v.Stat(path)
			}
		}
		if followSymlinks {
			return nil, fmt.Errorf("not supported")
		}
		return ws.FileInfo(path)
	}
}

func execEnv(env expand.Environ) []string {
	list := make([]string, 0, 64)
	for name, vr := range env.Each {
		if !vr.IsSet() {
			// If a variable is set globally but unset in the
			// runner, we need to ensure it's not part of the final
			// list. Seems like zeroing the element is enough.
			// This is a linear search, but this scenario should be
			// rare, and the number of variables shouldn't be large.
			for i, kv := range list {
				if strings.HasPrefix(kv, name+"=") {
					list[i] = ""
				}
			}
		}
		if vr.Exported && vr.Kind == expand.String {
			list = append(list, name+"="+vr.String())
		}
	}
	return list
}

func VirtualExecHandler(vs *VirtualSystem) func(next interp.ExecHandlerFunc) interp.ExecHandlerFunc {
	var killTimeout = 15 * time.Minute
	if vs.MaxTimeout > 0 {
		killTimeout = time.Duration(vs.MaxTimeout)
	}
	handle := func(ctx context.Context, args []string) error {
		hc := interp.HandlerCtx(ctx)
		path, err := interp.LookPathDir(hc.Dir, hc.Env, args[0])
		if err != nil {
			fmt.Fprintln(hc.Stderr, err)
			return interp.ExitStatus(127)
		}

		// cmd := exec.Cmd{
		// 	Path:   path,
		// 	Args:   args,
		// 	Env:    execEnv(hc.Env),
		// 	Dir:    hc.Dir,
		// 	Stdin:  hc.Stdin,
		// 	Stdout: hc.Stdout,
		// 	Stderr: hc.Stderr,
		// }

		// cmd := vs.System.Command(args[0], args[1:]...)
		cmd := vs.System.Command(path)
		cmd.Path = path
		cmd.Args = args
		cmd.Env = execEnv(hc.Env)
		cmd.Dir = hc.Dir
		cmd.Stdin = hc.Stdin
		cmd.Stdout = hc.Stdout
		cmd.Stderr = hc.Stderr

		prepareCommand(cmd)

		err = cmd.Start()
		if err == nil {
			stopf := context.AfterFunc(ctx, func() {
				if killTimeout <= 0 || runtime.GOOS == "windows" {
					_ = killCommand(cmd)
					return
				}
				_ = interruptCommand(cmd)
				// TODO: don't sleep in this goroutine if the program
				// stops itself with the interrupt above.
				time.Sleep(killTimeout)
				_ = killCommand(cmd)
			})
			defer stopf()

			err = cmd.Wait()
		}

		switch err := err.(type) {
		case *exec.ExitError:
			// Windows and Plan9 do not have support for [syscall.WaitStatus]
			// with methods like Signaled and Signal, so for those, [waitStatus] is a no-op.
			// Note: [waitStatus] is an alias [syscall.WaitStatus]
			if status, ok := err.Sys().(waitStatus); ok && status.Signaled() {
				if ctx.Err() != nil {
					return ctx.Err()
				}
				return interp.ExitStatus(128 + status.Signal())
			}
			return interp.ExitStatus(err.ExitCode())
		case *exec.Error:
			// did not start
			fmt.Fprintf(hc.Stderr, "%v\n", err)
			return interp.ExitStatus(127)
		default:
			return err
		}
	}

	return func(next interp.ExecHandlerFunc) interp.ExecHandlerFunc {
		return func(ctx context.Context, args []string) error {
			if err := handle(ctx, args); err != nil {
				return err
			}
			return next(ctx, args)
		}
	}
}

func NewRunner(vs *VirtualSystem, opts ...interp.RunnerOption) (*interp.Runner, error) {
	r, err := interp.New(opts...)
	if err != nil {
		return nil, err
	}

	interp.OpenHandler(VirtualOpenHandler(vs.Workspace))(r)
	interp.ReadDirHandler2(VirtualReadDirHandler2(vs.Workspace))(r)
	interp.StatHandler(VirtualStatHandler(vs.Workspace))(r)

	//
	var env = vs.System.Env()
	if len(env) > 0 {
		interp.Env(expand.ListEnviron(env...))(r)
	}

	dir, err := vs.System.Getwd()
	if err != nil {
		return nil, err
	}
	if err := interp.Dir(dir)(r); err != nil {
		return nil, err
	}
	interp.StdIO(vs.IOE.Stdin, vs.IOE.Stdout, vs.IOE.Stderr)(r)

	// exec handlers
	wrap := func(next interp.ExecHandlerFunc) interp.ExecHandlerFunc {
		return func(ctx context.Context, args []string) error {
			if vs.ExecHandler != nil {
				done, err := vs.ExecHandler(ctx, args)
				if done {
					return nil
				}
				if err != nil {
					return err
				}
			}
			return next(ctx, args)
		}
	}
	var middlewares = []func(interp.ExecHandlerFunc) interp.ExecHandlerFunc{
		// custom handler
		wrap,
		// default bash handler
		VirtualExecHandler(vs),
	}
	if err := interp.ExecHandlers(middlewares...)(r); err != nil {
		return nil, err
	}
	return r, nil
}
