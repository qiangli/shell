package sh

import (
	"context"
	"os"
	"testing"
)

func TestNewLocalSystem(t *testing.T) {
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
	}

	vs := NewLocalSystem("../", ioe)
	vs.ExecHandler = NewDummyExecHandler(vs)

	ctx := context.TODO()
	for _, tc := range tests {
		err := vs.RunScript(ctx, tc.script)
		if err != nil {
			t.FailNow()
		}
	}
}
