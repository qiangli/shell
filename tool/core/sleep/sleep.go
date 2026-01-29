// Copyright 2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Delay for the specified amount of time.
//
// Synopsis:
//
//	sleep DURATION
//
// Description:
//
//	If no units are given, the duration is assumed to be measured in
//	seconds, otherwise any format parsed by Go's `time.ParseDuration` is
//	accepted.
//
// Examples:
//
//	sleep 2.5
//	sleep 300ms
//	sleep 2h45m
//
// Bugs:
//
//	When sleep is first run, it must be compiled from source which creates a
//	delay significantly longer than anticipated.
package sleep

import (
	"context"
	"errors"
	"flag"
	"fmt"
	// "log"
	"time"

	"github.com/u-root/u-root/pkg/core"
	"github.com/u-root/u-root/pkg/uroot/unixflag"
)

var errDuration = errors.New("invalid duration")

func parseDuration(s string) (time.Duration, error) {
	d, err := time.ParseDuration(s)
	if err != nil {
		d, err = time.ParseDuration(s + "s")
	}
	if err != nil || d < 0 {
		return time.Duration(0), errDuration
	}
	return d, nil
}

// func main() {
// 	flag.Parse()

// 	if flag.NArg() != 1 {
// 		log.Fatal("Incorrect number of arguments")
// 	}

// 	d, err := parseDuration(flag.Arg(0))
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	time.Sleep(d)
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
}

// Run executes the command with a `context.Background()`.
func (c *command) Run(args ...string) error {
	return c.RunContext(context.Background(), args...)
}

// RunContext executes the command.
func (c *command) RunContext(ctx context.Context, args ...string) error {
	// var f flags

	fs := flag.NewFlagSet("sleep", flag.ContinueOnError)
	fs.SetOutput(c.Stderr)

	fs.Usage = func() {
		fmt.Fprintf(fs.Output(), "sleep DURATION\n\n")
		fmt.Fprintf(fs.Output(), "If no units are given, the duration is assumed to be measured in\n")
		fmt.Fprintf(fs.Output(), "seconds, otherwise any format parsed by Go's `time.ParseDuration` is\n")
		fmt.Fprintf(fs.Output(), "accepted.\n")
		fs.PrintDefaults()
	}

	if err := fs.Parse(unixflag.ArgsToGoArgs(args)); err != nil {
		return err
	}

	if fs.NArg() != 1 {
		// log.Fatal("Incorrect number of arguments")
		return fmt.Errorf("Incorrect number of arguments")
	}

	d, err := parseDuration(fs.Arg(0))
	if err != nil {
		// log.Fatal(err)
		return err
	}

	time.Sleep(d)

	return nil
}
