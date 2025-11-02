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
func Gosh(ctx context.Context, vs *VirtualSystem, args []string) error {
	script, args := parseFlags(args)
	err := runAll(ctx, vs, script, args)
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

func runAll(parent context.Context, vs *VirtualSystem, script string, args []string) error {
	ctx, _ := signal.NotifyContext(parent, os.Interrupt, syscall.SIGTERM)

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

	if file, ok := vs.IOE.Stdin.(*os.File); ok {
		if term.IsTerminal(int(file.Fd())) {
			return vs.Interactive(ctx)
		}
	}

	// piped
	return vs.RunStdin(ctx)
}

// Return script, non flag args
func parseFlags(args []string) (string, []string) {
	fs := flag.NewFlagSet("goshFlags", flag.ContinueOnError)
	var scriptptr = fs.String("c", "", "script to be executed")

	err := fs.Parse(args)
	if err != nil {
		fmt.Println("Error parsing flags:", err)
		return "", args
	}

	return *scriptptr, fs.Args()
}
