// Copyright 2021-2025 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package backoff

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"
)

func TestRunIt(t *testing.T) {

	for _, tt := range []struct {
		name    string
		timeout time.Duration
		cmd     string
		args    []string
		wantErr error
	}{
		{
			name:    "_date",
			timeout: 3 * time.Second,
			cmd:     "date",
		},
		{
			name:    "noCmd",
			timeout: 3 * time.Second,
			cmd:     "",
			wantErr: errNoCmd,
		},
		{
			name:    "echo",
			timeout: 3 * time.Second,
			cmd:     "echo",
			args:    []string{"hi"},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			run := func(to time.Duration, c string, a ...string) error {
				cb := func(args []string) error {
					cmd := exec.Command(args[0], args[1:]...)
					cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
					err := cmd.Run()
					return err
				}
				args := append([]string{"backoff", "-v", "-t", fmt.Sprintf("%v", to.Seconds()), c}, a...)
				t.Logf("testing %+v", args)
				bkoff := New(cb)
				return bkoff.Run(args[1:]...)
			}

			if err := run(tt.timeout, tt.cmd, tt.args...); !errors.Is(err, tt.wantErr) {
				if err != nil {
					if !strings.Contains(err.Error(), tt.wantErr.Error()) {
						t.Errorf("runit(%s, %s, %s)= %q, want %q", tt.timeout, tt.cmd, tt.args, err, tt.wantErr)
					}
				}
			}
		})
	}
}

// // TestOK for now just runs a simple successful test with 0 args or more than one arg.
// func TestOK(t *testing.T) {
// 	tests := []struct {
// 		args   []string
// 		stdout string
// 		stderr string
// 		exitok bool
// 	}{
// 		{args: []string{}, stdout: "", exitok: false},
// 		{args: []string{"date"}, stdout: ".*", exitok: true},
// 		{args: []string{"echo", "hi"}, stdout: ".*hi", exitok: true},
// 		{args: []string{"-t", "1s", "false"}, exitok: false},
// 	}

// 	for _, v := range tests {
// 		c := testutil.Command(t, v.args...)
// 		stdout, stderr := &bytes.Buffer{}, &bytes.Buffer{}
// 		c.Stdout, c.Stderr = stdout, stderr
// 		err := c.Run()
// 		if (err != nil) && v.exitok {
// 			t.Errorf("%v: got %v, want nil", v, err)
// 		}
// 		if (err == nil) && !v.exitok {
// 			t.Errorf("%v: got nil, want err", v)
// 		}
// 		m, err := regexp.MatchString(v.stderr, stderr.String())
// 		if err != nil {
// 			t.Errorf("stderr: %v: got %v, want nil", v, err)
// 		} else {
// 			if !m {
// 				t.Errorf("%v: regexp.MatchString(%s, %s) false, wanted match", v, v.stderr, stderr)
// 			}
// 		}

// 		m, err = regexp.MatchString(v.stdout, stdout.String())
// 		if err != nil {
// 			t.Errorf("stdout: %v: got %v, want nil", v, err)
// 		}
// 		if !m {
// 			t.Errorf("%v: regexp.MatchString(%s, %s) false, wanted match", v, v.stdout, stderr.String())
// 		}
// 	}
// }

// // If you really like fork-bombing your machine, remove these lines :-)
// func TestMain(m *testing.M) {
// 	testutil.Run(m, main)
// }
