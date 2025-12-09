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
	ws, script, args := sh.ParseFlags(os.Args[1:])
	if ws == "" {
		ws, _ = os.Getwd()
	}

	ws, _ = filepath.Abs(ws)
	if err := os.Chdir(ws); err != nil {
		fmt.Printf("%v", err)
		os.Exit(1)
	}

	home, _ := os.UserHomeDir()
	tmpdir := os.TempDir()

	lfs, _ := vfs.NewLocalFS([]string{ws, home, tmpdir})
	los, _ := vos.NewLocalSystem(lfs)
	los.Exitf = func(code int) {
		fmt.Printf("exit %v\n", code)
		os.Exit(code)
	}

	ioe := &sh.IOE{Stdin: os.Stdin, Stdout: os.Stdout, Stderr: os.Stderr}
	vs := sh.NewVirtualSystem(los, lfs, ioe)
	vs.ExecHandler = sh.NewDummyExecHandler(vs)

	if err := sh.Gosh(context.Background(), vs, script, args); err != nil {
		os.Exit(1)
	}
}
