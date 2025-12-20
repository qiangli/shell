// Copyright 2012-2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Time process execution.
//
// Synopsis:
//
//	time [-p] CMD [ARG]...
//
// Description:
//
//	After executing CMD, its real, user and system times are printed to
//	stderr in the POSIX format.
//
// Example:
//
//	$ time sleep 1.23s
//	real 1.230
//	user 0.001
//	sys 0.000
//
// Note:
//
//	This is different from bash's time command which is built into the shell
//	and can time the entire pipeline.
//
// Bugs:
//
//	Time is not reported when exiting due to a signal.
package time

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/u-root/u-root/pkg/core"
	"github.com/u-root/u-root/pkg/uroot/unixflag"
)

// var _ = flag.Bool("p", true, "makes time output POSIX.2 compliant")

func printTime(stderr io.Writer, l string, t time.Duration) {
	fmt.Fprintf(stderr, "%s\n", label(l, t))
}

func label(l string, t time.Duration) string {
	return fmt.Sprintf("%s %.03f", l, t.Seconds())
}

func run(args []string, stdin io.Reader, stdout io.Writer, stderr io.Writer) error {
	start := time.Now()
	if len(args) == 0 {
		fmt.Fprintf(stderr, "%s\n", strings.Join([]string{
			label("real", 0*time.Second),
			label("user", 0*time.Second),
			label("sys", 0*time.Second),
		}, "\n"))
		return nil
	}
	c := exec.Command(args[0], args[1:]...)
	c.Stdin, c.Stdout, c.Stderr = stdin, stdout, stderr
	defer func(*exec.Cmd, time.Time) {
		realTime := time.Since(start)
		printTime(stderr, "real", realTime)
		printProcessState(stderr, c)
	}(c, start)
	if err := c.Run(); err != nil {
		return fmt.Errorf("%q:%w", args, err)
	}
	return nil
}

// func main() {
// 	flag.Parse()
// 	if err := run(flag.Args(), os.Stdin, os.Stdout, os.Stderr); err != nil {
// 		log.Fatalf("time: %v", err)
// 	}
// }

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
	posix bool
}

// Run executes the command with a `context.Background()`.
func (c *command) Run(args ...string) error {
	return c.RunContext(context.Background(), args...)
}

// RunContext executes the command.
func (c *command) RunContext(ctx context.Context, args ...string) error {
	var f flags

	fs := flag.NewFlagSet("time", flag.ContinueOnError)
	fs.SetOutput(c.Stderr)

	fs.BoolVar(&f.posix, "p", true, "makes time output POSIX.2 compliant")

	fs.Usage = func() {
		fmt.Fprintf(fs.Output(), "time [-p] CMD [ARG]...\n\n")
		fmt.Fprintf(fs.Output(), "After executing CMD, its real, user and system times are printed to\n")
		fmt.Fprintf(fs.Output(), "stderr in the POSIX format.\n")
		fs.PrintDefaults()
	}

	if err := fs.Parse(unixflag.ArgsToGoArgs(args)); err != nil {
		return err
	}

	// if err := run(flag.Args(), f.posix, rc, c.Stdout); err != nil {
	// 	log.Fatalf("date: %v", err)
	// }

	if err := run(flag.Args(), os.Stdin, os.Stdout, os.Stderr); err != nil {
		log.Fatalf("time: %v", err)
	}

	return nil
}
