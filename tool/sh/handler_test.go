package sh

import (
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

	ls := NewLocalSystem("../", ioe)
	ls.ExecHandler = NewDummyExecHandler(ioe)

	for _, tc := range tests {
		err := ls.Run(tc.script)
		if err != nil {
			t.FailNow()
		}
	}
}
