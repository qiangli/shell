// Copyright 2016 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dirname

import (
	"bytes"
	"errors"
	"testing"
)

func TestDirName(t *testing.T) {

	type test struct {
		args []string
		out  string
		err  error
	}

	var dirnameTests = []test{
		// For no args it seems we have to print an error.
		// It should be missing operand[s] but that's not the standard.
		{args: []string{}, err: ErrNoArg},
		{args: []string{""}, out: ".\n"},
		{args: []string{"/this/that"}, out: "/this\n"},
		{args: []string{"/this/that", "/other"}, out: "/this\n/\n"},
		{args: []string{"/this/that", "/other thing/space"}, out: "/this\n/other thing\n"},
	}

	// Table-driven testing
	out := bytes.NewBuffer(nil)

	for _, tt := range dirnameTests {
		out.Reset()

		cmd := New()
		cmd.SetIO(nil, out, out)
		args := tt.args
		err := cmd.Run(args...)

		if !errors.Is(err, tt.err) {
			t.Errorf("errors do not match: got %v - want %v", err, tt.err)
		}

		if out.String() != tt.out {
			t.Errorf("%v: got %q, wants %q", tt.args, out.String(), tt.out)
		}

	}
}
