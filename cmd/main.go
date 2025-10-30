package main

import (
	"os"

	"github.com/qiangli/shell/tool/sh"
)

func main() {
	ioe := &sh.IOE{Stdin: os.Stdin, Stdout: os.Stdout, Stderr: os.Stderr}
	args := os.Args[1:]
	sh.Gosh(sh.NewLocalSystem(ioe), args)
}
