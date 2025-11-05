// Copyright 2012-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// tac concatenates files and prints to stdout in reverse order,
// file by file
//
// Synopsis:
//
//	tac <file...>
//
// Description:
//
// Options:
package tac

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"sync"

	"github.com/u-root/u-root/pkg/core"
	"github.com/u-root/u-root/pkg/uroot/unixflag"
)

const ReadSize int64 = 4096

var errStdin = fmt.Errorf("can't reverse lines from stdin; can't seek")

type ReadAtSeeker interface {
	io.ReaderAt
	io.Seeker
}

func tacOne(w io.Writer, r ReadAtSeeker) error {
	var b [ReadSize]byte
	// Get current EOF. While the file may be growing, there's
	// only so much we can do.
	loc, err := r.Seek(0, io.SeekEnd)
	if err != nil {
		return err
	}
	var wg sync.WaitGroup
	wg.Add(1)
	c := make(chan byte)
	go func(r <-chan byte, w io.Writer) {
		defer wg.Done()
		line := string(<-r)
		for c := range r {
			if c == '\n' {
				if _, err := w.Write([]byte(line)); err != nil {
					log.Fatal(err)
				}
				line = ""
			}
			line = string(c) + line
		}
		if _, err := w.Write([]byte(line)); err != nil {
			log.Fatal(err)
		}
	}(c, w)

	for loc > 0 {
		n := min(loc, ReadSize)

		amt, err := r.ReadAt(b[:n], loc-int64(n))
		if err != nil && err != io.EOF {
			return err
		}
		loc -= int64(amt)
		for i := range b[:amt] {
			o := amt - i - 1
			c <- b[o]
		}
	}
	close(c)
	wg.Wait()
	return nil
}

func tac(w io.Writer, files []string) error {
	if len(files) == 0 {
		return errStdin
	}
	for _, name := range files {
		f, err := os.Open(name)
		if err != nil {
			return err
		}
		err = tacOne(w, f)
		f.Close() // Don't defer, you might get EMFILE for no good reason.
		if err != nil {
			return err
		}

	}
	return nil
}

// func main() {
// 	flag.Parse()
// 	if err := tac(os.Stdout, flag.Args()); err != nil {
// 		log.Fatalf("tac: %v", err)
// 	}
// }

// command implements the cat core utility.
type command struct {
	core.Base
}

// New creates a new tac command.
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
	fs := flag.NewFlagSet("tac", flag.ContinueOnError)
	fs.SetOutput(c.Stderr)

	fs.Usage = func() {
		fmt.Fprintf(fs.Output(), "Usage: tac <file...>\n\n")
		fmt.Fprintf(fs.Output(), "tac concatenates files and prints to stdout in reverse order\n")
		fmt.Fprintf(fs.Output(), "file by file.\n\n")
		fmt.Fprintf(fs.Output(), "Options:\n")
		fs.PrintDefaults()
	}

	if err := fs.Parse(unixflag.ArgsToGoArgs(args)); err != nil {
		return err
	}

	if err := tac(c.Stdout, fs.Args()); err != nil {
		return err
	}

	return nil
}
