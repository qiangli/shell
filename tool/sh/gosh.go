package sh

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"golang.org/x/term"

	"mvdan.cc/sh/v3/interp"
	"mvdan.cc/sh/v3/syntax"
)

// Gosh executes a script provided in the argument.
// If no arguments are provided, it will execute in interactive mode
// if standard input supports it.
// This function manages errors and exits appropriately.
func Gosh(vs *VirtualSystem, script string) {
	err := runAll(vs, script)
	var es interp.ExitStatus
	if errors.As(err, &es) {
		vs.System.Exit(int(es))
	}
	if err != nil {
		fmt.Fprintln(vs.IOE.Stderr, err)
		vs.System.Exit(1)
	}
}

func runAll(vs *VirtualSystem, script string) error {
	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)

	r, err := NewRunner(vs, interp.Interactive(true))
	if err != nil {
		return err
	}

	if script != "" {
		return run(ctx, r, strings.NewReader(script), "")
	}

	if file, ok := vs.IOE.Stdin.(*os.File); ok {
		if term.IsTerminal(int(file.Fd())) {
			return runInteractive(vs, ctx, r)
		}
	}
	return run(ctx, r, vs.IOE.Stdin, "")
}

func run(ctx context.Context, r *interp.Runner, reader io.Reader, name string) error {
	prog, err := syntax.NewParser().Parse(reader, name)
	if err != nil {
		return err
	}
	r.Reset()
	return r.Run(ctx, prog)
}

func runInteractive(vs *VirtualSystem, ctx context.Context, r *interp.Runner) error {
	parser := syntax.NewParser()

	fmt.Fprintf(vs.IOE.Stdout, "$ ")
	err := parser.Interactive(vs.IOE.Stdin, func(stmts []*syntax.Stmt) bool {
		if parser.Incomplete() {
			fmt.Fprintf(vs.IOE.Stdout, "> ")
			return true
		}
		// run
		for _, stmt := range stmts {
			err := r.Run(ctx, stmt)
			if err != nil {
				fmt.Fprint(vs.IOE.Stderr, err.Error())
			}
			if r.Exited() {
				vs.System.Exit(0)
				return true
			}
		}
		fmt.Fprintf(vs.IOE.Stdout, "$ ")
		return true
	})
	return err
}
