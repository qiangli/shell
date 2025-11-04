// Copyright 2016-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// dirname prints out the directory name of one or more args.
// If no arg is given it returns an error and prints a message which,
// per the man page, is incorrect, but per the standard, is correct.
package dirname

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"path/filepath"

	"github.com/u-root/u-root/pkg/core"
	"github.com/u-root/u-root/pkg/uroot/unixflag"
)

var ErrNoArg = errors.New("missing operand")

func (c *command) run(out io.Writer, args []string) error {
	if len(args) < 1 {
		return ErrNoArg
	}

	for _, n := range args {
		fmt.Fprintln(out, filepath.Dir(n))
	}
	return nil
}

// func main() {
// 	if err := run(os.Stdout, os.Args[1:]); err != nil {
// 		log.Fatalf("dirname: %v", err)
// 	}
// }

// command implements the cat core utility.
type command struct {
	core.Base
}

// New creates a new cat command.
func New() core.Command {
	c := &command{}
	c.Init()
	return c
}

type flags struct {
}

// Run executes the command with a `context.Background()`.
func (c *command) Run(args ...string) error {
	return c.RunContext(context.Background(), args...)
}

// RunContext executes the command.
func (c *command) RunContext(ctx context.Context, args ...string) error {
	// var f flags

	fs := flag.NewFlagSet("cat", flag.ContinueOnError)
	fs.SetOutput(c.Stderr)

	// fs.BoolVar(&f.u, "u", false, "ignored")

	fs.Usage = func() {
		fmt.Fprintf(fs.Output(), "dirname NAME [SUFFIX]\n\n")
		fmt.Fprintf(fs.Output(), "dirname return directory portion of pathname.\n")
		fmt.Fprintf(fs.Output(), "Options:\n")
		fs.PrintDefaults()
	}

	if err := fs.Parse(unixflag.ArgsToGoArgs(args)); err != nil {
		return err
	}

	if err := c.run(c.Stdout, fs.Args()); err != nil {
		return err
	}

	return nil
}
