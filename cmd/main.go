package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/qiangli/shell/tool/sh"
	"github.com/qiangli/shell/tool/sh/vfs"
	"github.com/qiangli/shell/tool/sh/vos"
)

func main() {
	root, script, args := sh.ParseFlags(os.Args[1:])
	if root == "" {
		root, _ = os.Getwd()
	}

	root, _ = filepath.Abs(root)
	if err := os.Chdir(root); err != nil {
		fmt.Printf("%v", err)
		os.Exit(1)
	}

	los := vos.NewLocalSystem(root)
	los.Exitf = func(code int) {
		fmt.Printf("exit %v\n", code)
		os.Exit(code)
	}

	lfs := vfs.NewLocalFS(root)
	ioe := &sh.IOE{Stdin: os.Stdin, Stdout: os.Stdout, Stderr: os.Stderr}
	vs := sh.NewVirtualSystem(root, los, lfs, ioe)
	vs.ExecHandler = sh.NewDummyExecHandler(vs)

	if err := sh.Gosh(context.Background(), vs, script, args); err != nil {
		os.Exit(1)
	}
}
