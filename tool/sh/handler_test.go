package sh

import (
	"context"
	_ "embed"
	"os"
	"testing"
)

//go:embed testdata/test.sh
var test_sh string

func TestRunScript(t *testing.T) {
	ioe := &IOE{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
	tests := []struct {
		script string
	}{
		{"echo hello world"},
		{"ai --models hello"},
		{"@agent --max-history 0 hello"},
		{"@ anonymous agent"},
		{test_sh},
	}

	vs := NewLocalSystem("./", ioe)
	vs.ExecHandler = NewDummyExecHandler(vs)

	ctx := context.TODO()
	for _, tc := range tests {
		err := vs.RunScript(ctx, tc.script)
		if err != nil {
			t.FailNow()
		}
	}
}

func TestRunPath(t *testing.T) {
	ioe := &IOE{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
	tests := []struct {
		script string
	}{
		// {"testdata/test.sh"},
		{"testdata/coreutils.sh"},
	}

	vs := NewLocalSystem("./", ioe)
	vs.ExecHandler = NewDummyExecHandler(vs)

	vs.System.Setenv("city", "New York")

	ctx := context.TODO()
	for _, tc := range tests {
		err := vs.RunPath(ctx, tc.script)
		if err != nil {
			t.FailNow()
		}
	}
}
