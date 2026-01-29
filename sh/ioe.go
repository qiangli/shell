package sh

import (
	"io"
	"os"
)

// standard IO
type IOE struct {
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

// type IOE interface {
// 	Stdin()  io.Reader
// 	Stdout() io.Writer
// 	Stderr() io.Writer
// }

// type StdIOE struct {
// }

// func (r StdIOE) Stdin() io.Reader {
// 	return os.Stdin
// }

// func (r StdIOE) Stdout() io.Writer {
// 	return os.Stdout
// }

// func (r StdIOE) Stderr() io.Writer {
// 	return os.Stderr
// }

func NewStdIOE() *IOE {
	return &IOE{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
}

type StringIOE struct {
	stdin  io.Reader
	stdout io.Writer
	stderr io.Writer
}

func NewStringIOE(s string) *IOE {

	return &IOE{}
}
