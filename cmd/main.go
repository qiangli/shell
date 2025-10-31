package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/qiangli/shell/tool/sh"
	"github.com/qiangli/shell/tool/sh/vfs"
	"github.com/qiangli/shell/tool/sh/vos"
)

func main() {
	var script = flag.String("c", "", "script to be executed")
	var root = flag.String("root", "", "Specify the root directory")
	flag.Parse()

	// args := flag.Args()
	if *root != "" {
		if err := os.Chdir(*root); err != nil {
			fmt.Printf("%v", err)
			os.Exit(1)
		}
	}

	ls := vos.NewLocalSystem(*root)
	ls.Exitf = func(code int) {
		fmt.Printf("exit %v\n", code)
		os.Exit(code)
	}
	ioe := &sh.IOE{Stdin: os.Stdin, Stdout: os.Stdout, Stderr: os.Stderr}
	vs := sh.NewVirtualSystem(ls, vfs.NewLocalFS(), ioe)
	sh.Gosh(vs, *script)
}
