package main

import (
	"os"

	"github.com/qiangli/shell/tool/sh"
)

func main() {
	ioe := [3]*os.File{os.Stdin, os.Stdout, os.Stderr}
	args := os.Args[1:]
	sh.Gosh(sh.NewLocalSystem(ioe), args)
}
