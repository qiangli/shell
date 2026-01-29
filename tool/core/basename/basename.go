// Copyright 2012-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Basename return name with leading path information removed.
//
// Synopsis:
//
//	basename NAME [SUFFIX]
package basename

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/u-root/u-root/pkg/core"
	"github.com/u-root/u-root/pkg/uroot/unixflag"
)

var errUsage = errors.New("usage: basename NAME [SUFFIX]")

func run(w io.Writer, args []string) error {
	switch len(args) {
	case 2:
		fileName := filepath.Base(args[0])
		if fileName != args[1] {
			fileName = strings.TrimSuffix(fileName, args[1])
		}
		_, err := fmt.Fprintf(w, "%s\n", fileName)
		return err
	case 1:
		fileName := filepath.Base(args[0])
		_, err := fmt.Fprintf(w, "%s\n", fileName)
		return err
	default:
		return errUsage
	}
}

// func main() {
// 	if err := run(os.Stdout, os.Args[1:]); err != nil {
// 		log.Fatal(err)
// 	}
// }

// command implements the basename core utility.
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

	fs := flag.NewFlagSet("basename", flag.ContinueOnError)
	fs.SetOutput(c.Stderr)

	// fs.BoolVar(&f.u, "u", false, "ignored")

	fs.Usage = func() {
		fmt.Fprintf(fs.Output(), "basename NAME [SUFFIX]\n\n")
		fmt.Fprintf(fs.Output(), "Basename return name with leading path information removed.\n")
		fmt.Fprintf(fs.Output(), "Options:\n")
		fs.PrintDefaults()
	}

	if err := fs.Parse(unixflag.ArgsToGoArgs(args)); err != nil {
		return err
	}

	if err := run(c.Stdout, fs.Args()); err != nil {
		return err
	}

	return nil
}
