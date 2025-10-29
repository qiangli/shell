package vos

import (
	// "bytes"
	// "fmt"
	"os"
	"os/exec"
	// "strings"
)

// System represents the virtual operating system for the tool.
// It provides the system operations that can be mocked for testing.
type System interface {
	// Man(string) (string, error)

	Command(name string, arg ...string) *exec.Cmd

	Chdir(dir string) error
	Getwd() (string, error)

	Environ() []string
	Getenv(string) string
	Setenv(string, string)

	Exit(int)
}

type LocalSystem struct {
}

func NewLocalSystem() *LocalSystem {
	return &LocalSystem{}
}

func (s *LocalSystem) Command(name string, arg ...string) *exec.Cmd {
	return exec.Command(name, arg...)
}

func (s *LocalSystem) Chdir(dir string) error {
	return os.Chdir(dir)
}

func (s *LocalSystem) Getwd() (string, error) {
	return os.Getwd()
}

// func (s *LocalSystem) Man(bin string) (string, error) {
// 	command := strings.TrimSpace(strings.SplitN(bin, " ", 2)[0])
// 	manCmd := s.Command("man", command)
// 	var manOutput bytes.Buffer

// 	// Capture the output of the man command.
// 	manCmd.Stdout = &manOutput
// 	manCmd.Stderr = &manOutput

// 	if err := manCmd.Run(); err != nil {
// 		return "", fmt.Errorf("error running man command: %v\nOutput: %s", err, manOutput.String())
// 	}

// 	// Process with 'col' to remove formatting
// 	colCmd := s.Command("col", "-b")
// 	var colOutput bytes.Buffer

// 	colCmd.Stdin = bytes.NewReader(manOutput.Bytes())
// 	colCmd.Stdout = &colOutput
// 	colCmd.Stderr = &colOutput

// 	// Try running 'col', if it fails, return the man output instead.
// 	if err := colCmd.Run(); err != nil {
// 		return manOutput.String(), nil
// 	}

// 	return colOutput.String(), nil
// }

func (s *LocalSystem) Environ() []string {
	return os.Environ()
}

func (s *LocalSystem) Getenv(key string) string {
	return os.Getenv(key)
}

func (s *LocalSystem) Setenv(key string, value string) {
	os.Setenv(key, value)
}

func (s *LocalSystem) Exit(code int) {
	os.Exit(code)
}
