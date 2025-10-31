package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/qiangli/shell/tool/sh"
	"github.com/qiangli/shell/tool/sh/vfs"
	"github.com/qiangli/shell/tool/sh/vos"
)

func main() {
	var scriptptr = flag.String("c", "", "script to be executed")
	var rootptr = flag.String("root", "", "Specify the root directory")
	flag.Parse()

	var script = *scriptptr
	var root = *rootptr

	if root == "" {
		root, _ = os.Getwd()
	}

	root, _ = filepath.Abs(root)
	if err := os.Chdir(root); err != nil {
		fmt.Printf("%v", err)
		os.Exit(1)
	}

	ls := vos.NewLocalSystem(root)
	ls.Exitf = func(code int) {
		fmt.Printf("exit %v\n", code)
		os.Exit(code)
	}
	ioe := &sh.IOE{Stdin: os.Stdin, Stdout: os.Stdout, Stderr: os.Stderr}
	vs := sh.NewVirtualSystem(ls, vfs.NewLocalFS(root), ioe)
	sh.Gosh(vs, script)
}
