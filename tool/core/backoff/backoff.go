// Copyright 2012-2025 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Run a command, repeatedly, until it succeeds or we are out of time
//
// Synopsis:
//	backoff [-v] [-t duration-string] command [args...]
//
// Description:
//	backoff will run the command until it succeeds or a timeout has passed.
//	The default timeout is 30s.
//	If -v is set, it will show what it is running, each time it is tried.
//	If no args are given, it will print command help.
//
// Example:
//	$ backoff echo hi
//	hi
//	$
//	$ backoff -v -t=2s false
//	  2022/03/31 14:29:37 Run ["false"]
//	  2022/03/31 14:29:37 Set timeout to 2s
//	  2022/03/31 14:29:37 "false" []:exit status 1
//	  2022/03/31 14:29:38 "false" []:exit status 1
//	  2022/03/31 14:29:39 "false" []:exit status 1
//	  2022/03/31 14:29:39 Error: exit status 1

//go:build !test

package backoff

import (
	"context"
	"flag"
	"fmt"
	"time"

	"github.com/cenkalti/backoff/v4"

	"github.com/u-root/u-root/pkg/core"
	"github.com/u-root/u-root/pkg/uroot/unixflag"
)

var (
	errNoCmd = fmt.Errorf("no command passed")
)

func (c *command) run(_ context.Context, timeout time.Duration, verbose bool, args []string) error {
	if args[0] == "" {
		return errNoCmd
	}
	if verbose {
		fmt.Fprintf(c.Stdout, "Run %q", args)
	}

	b := backoff.NewExponentialBackOff()
	b.MaxElapsedTime = timeout
	f := func() error {
		err := c.cb(args)
		if verbose {
			fmt.Fprintf(c.Stdout, "Run %q: %v", args, err)
		}
		return err
	}

	return backoff.Retry(f, b)
}

// func main() {
// 	flag.Parse()
// 	if *verbose {
// 		v = log.Printf
// 	}
// 	a := flag.Args()
// 	if len(a) == 0 {
// 		flag.Usage()
// 		os.Exit(1)
// 	}
// 	v("Run %q", a)
// 	if err := run(*timeout, a[0], a[1:]...); err != nil {
// 		log.Fatalf("Error: %v", err)
// 	}
// }

type Callback func([]string) error

type command struct {
	core.Base

	cb Callback
}

// New creates a new command.
func New(cb Callback) core.Command {
	c := &command{
		cb: cb,
	}
	c.Init()
	return c
}

type flags struct {
	timeout int64
	verbose bool
}

// Run executes the command with a `context.Background()`.
func (c *command) Run(args ...string) error {
	return c.RunContext(context.Background(), args...)
}

// RunContext executes the command.
func (c *command) RunContext(ctx context.Context, args ...string) error {
	var f flags

	fs := flag.NewFlagSet("backoff", flag.ContinueOnError)

	fs.Int64Var(&f.timeout, "t", 30, "Timeout for command in seconds")
	fs.BoolVar(&f.verbose, "v", false, "Log each attempt to run the command")

	fs.SetOutput(c.Stderr)

	fs.Usage = func() {
		fmt.Fprintf(fs.Output(), "Usage: backoff [-v] [-t duration-string] command [args...]\n\n")
		fmt.Fprintf(fs.Output(), "backoff will run the command until it succeeds or a timeout has passed.\n")
		fmt.Fprintf(fs.Output(), "The default timeout is 30s.\n\n")
		fmt.Fprintf(fs.Output(), "Options:\n")
		fs.PrintDefaults()
	}

	if err := fs.Parse(unixflag.ArgsToGoArgs(args)); err != nil {
		return err
	}

	return c.run(ctx, time.Duration(f.timeout)*time.Second, f.verbose, fs.Args())
}
