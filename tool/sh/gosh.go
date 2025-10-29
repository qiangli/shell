package sh

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/term"

	"mvdan.cc/sh/v3/interp"
	"mvdan.cc/sh/v3/syntax"
)

func Gosh(vs *VirtualSystem, args []string) {
	err := runAll(vs, args)
	var es interp.ExitStatus
	if errors.As(err, &es) {
		vs.System.Exit(int(es))
	}
	if err != nil {
		fmt.Fprintln(vs.IOE[2], err)
		vs.System.Exit(1)
	}
}

func runAll(vs *VirtualSystem, args []string) error {
	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)

	r, err := NewRunner(vs, interp.Interactive(true))
	if err != nil {
		return err
	}

	if len(args) == 0 {
		if term.IsTerminal(int(vs.IOE[0].Fd())) {
			return runInteractive(vs, ctx, r)
		}
		return run(ctx, r, vs.IOE[0], "")
	}
	for _, path := range args {
		if err := runPath(vs, ctx, r, path); err != nil {
			return err
		}
	}
	return nil

}

func run(ctx context.Context, r *interp.Runner, reader io.Reader, name string) error {
	prog, err := syntax.NewParser().Parse(reader, name)
	if err != nil {
		return err
	}
	r.Reset()
	return r.Run(ctx, prog)
}

func runPath(vs *VirtualSystem, ctx context.Context, r *interp.Runner, path string) error {
	// f, err := os.Open(path)
	f, err := vs.Workspace.OpenFile(path, os.O_RDONLY, 0)
	if err != nil {
		return err
	}
	defer f.Close()
	return run(ctx, r, f, path)
}

func runInteractive(vs *VirtualSystem, ctx context.Context, r *interp.Runner) error {
	parser := syntax.NewParser()

	fmt.Fprintf(vs.IOE[1], "$ ")
	err := parser.Interactive(vs.IOE[0], func(stmts []*syntax.Stmt) bool {
		if parser.Incomplete() {
			fmt.Fprintf(vs.IOE[1], "> ")
			return true
		}
		// run
		for _, stmt := range stmts {
			err := r.Run(ctx, stmt)
			if err != nil {
				fmt.Fprint(vs.IOE[2], err.Error())
			}
			if r.Exited() {
				vs.System.Exit(0)
				return true
			}
		}
		fmt.Fprintf(vs.IOE[1], "$ ")
		return true
	})
	return err
}
