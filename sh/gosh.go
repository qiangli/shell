package sh

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/term"

	"mvdan.cc/sh/v3/interp"
)

// Gosh executes a script provided in the argument.
// If no arguments are provided, it will execute in interactive mode
// if standard input supports it.
// This function manages errors and exits appropriately.
func Gosh(ctx context.Context, vs *VirtualSystem, script string, args []string) error {
	err := Run(ctx, vs, script, args)
	var es interp.ExitStatus
	if errors.As(err, &es) {
		vs.System.Exit(int(es))
	}
	if err != nil {
		fmt.Fprintln(vs.IOE.Stderr, err)
		vs.System.Exit(1)
	}
	return err
}

func Run(parent context.Context, vs *VirtualSystem, script string, args []string) error {
	ctx, _ := signal.NotifyContext(parent, os.Interrupt, syscall.SIGTERM)
	// defer cancel()

	if script != "" {
		return vs.RunScript(ctx, script)
	}

	if len(args) > 0 {
		for _, path := range args {
			if err := vs.RunPath(ctx, path); err != nil {
				return err
			}
		}
		return nil
	}

	if v, ok := vs.IOE.Stdin.(*os.File); ok && term.IsTerminal(int(v.Fd())) {
		return vs.RunInteractive(ctx)
	}

	// piped
	return vs.RunReader(ctx)
}

// Parse parses flag definitions from the argument list, which should not
// include the command name.
// Return root, script, and remaining non flag args
func ParseFlags(args []string) (string, string, []string) {
	fs := flag.NewFlagSet("gosh", flag.ContinueOnError)
	var rootptr = fs.String("root", "", "Specify the workspace root directory")
	var command = fs.String("c", "", "script to run")

	err := fs.Parse(args)

	if err != nil {
		fmt.Println("Error parsing flags:", err)
		return "", "", args
	}

	return *rootptr, *command, fs.Args()
}
