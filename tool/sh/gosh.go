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
	ctx, cancel := signal.NotifyContext(parent, os.Interrupt, syscall.SIGTERM)
	defer cancel()

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

// Return script, non flag args
func parseFlags(args []string) (string, []string) {
	fs := flag.NewFlagSet("gosh", flag.ContinueOnError)
	var command = fs.String("c", "", "script to run")

	err := fs.Parse(args)
	if err != nil {
		fmt.Println("Error parsing flags:", err)
		return "", args
	}

	return *command, fs.Args()
}
