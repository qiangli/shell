// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package head

import (
	"bufio"
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"

	"github.com/u-root/u-root/pkg/core"
	"github.com/u-root/u-root/pkg/uroot/unixflag"
)

var errCombine = fmt.Errorf("can't combine line and byte counts")

func (c *command) run(stdin io.Reader, stdout, stderr io.Writer, bytes, count int, files ...string) error {
	if bytes > 0 && count > 0 {
		return errCombine
	}

	var printBytes bool
	var buffer []byte
	if bytes > 0 {
		printBytes = true
		buffer = make([]byte, 4096)
	}

	if count == 0 {
		count = 10
	}

	var newLineHeader bool
	var errs error

	handle := func(r io.Reader, name string) error {
		if len(files) > 1 {
			if newLineHeader {
				fmt.Fprintf(stdout, "\n==> %s <==\n", name)
			} else {
				fmt.Fprintf(stdout, "==> %s <==\n", name)
				newLineHeader = true
			}
		}
		if printBytes {
			c := bytes
			for {
				n, err := io.ReadFull(r, buffer)
				if err == io.EOF {
					break
				}
				if err != nil && err != io.ErrUnexpectedEOF {
					return err
				}

				stdout.Write(buffer[:min(c, n)])
				c -= n
				if c <= 0 {
					break
				}

				// handle the case when user request more bytes than
				// source have
				if err == io.ErrUnexpectedEOF {
					break
				}
			}
		} else {
			var c int
			scanner := bufio.NewScanner(r)
			for scanner.Scan() {
				fmt.Fprintln(stdout, scanner.Text())
				c++
				if c == count {
					break
				}
			}
		}
		return nil
	}

	// handle stdin
	if len(files) == 0 {
		return handle(stdin, "")
	}

	for _, file := range files {
		f, err := c.f.Open(file)
		if err != nil {
			errs = errors.Join(errs, fmt.Errorf("head: %w", err))
			continue
		}
		err = handle(f, file)
		if err != nil {
			errs = errors.Join(errs, err)
		}
	}

	if errs != nil {
		fmt.Fprintf(stderr, "\n%v\n", errs)
	}
	return nil
}

// func main() {
// 	c := flag.Int("c", 0, "Print bytes of each of the specified files")
// 	n := flag.Int("n", 0, "Print count lines of each of the specified files")

// 	flag.Parse()
// 	if err := run(os.Stdin, os.Stdout, os.Stderr, *c, *n, flag.Args()...); err != nil {
// 		log.Fatalf("head: %v", err)
// 	}
// }

// command implements the head core utility.
// type FileOpen func(string) (*os.File, error)
type command struct {
	core.Base

	f fs.FS
}

// New creates a new cat command.
func New(f fs.FS) core.Command {
	c := &command{
		f: f,
	}
	c.Init()
	return c
}

type flags struct {
	c int
	n int
}

// Run executes the command with a `context.Background()`.
func (c *command) Run(args ...string) error {
	return c.RunContext(context.Background(), args...)
}

// RunContext executes the command.
func (c *command) RunContext(ctx context.Context, args ...string) error {
	var f flags

	fs := flag.NewFlagSet("head", flag.ContinueOnError)
	fs.SetOutput(c.Stderr)

	fs.IntVar(&f.c, "c", 0, "Print bytes of each of the specified files")
	fs.IntVar(&f.n, "n", 0, "Print count lines of each of the specified files")

	fs.Usage = func() {
		fmt.Fprintf(fs.Output(), "head [-n count | -c bytes] [file ...]\n\n")
		fmt.Fprintf(fs.Output(), "head -- display first lines of a file.\n")
		fmt.Fprintf(fs.Output(), "This filter displays the first count lines or bytes of each of the specified files,\n")
		// fmt.Fprintf(fs.Output(), "or of the standard input if no files are specified.  If count is omitted it defaults\n")
		// fmt.Fprintf(fs.Output(), "to 10.\n")
		fmt.Fprintf(fs.Output(), "Options:\n")
		fs.PrintDefaults()
	}

	if err := fs.Parse(unixflag.ArgsToGoArgs(args)); err != nil {
		return err
	}

	if err := c.run(c.Stdin, c.Stdout, c.Stderr, f.c, f.n, fs.Args()...); err != nil {
		return err
	}
	return nil
}
