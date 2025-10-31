package vos

import (
	"os/exec"
)

type System interface {
	Command(name string, arg ...string) *exec.Cmd
	Chdir(dir string) error
	Getwd() (string, error)
	Env() []string
	Environ() map[string]any
	Getenv(string) any
	Setenv(string, any)
	Exit(int)
}
