package main

import (
	"context"
	"fmt"
	"os"

	"github.com/zdz1715/go-sh"
)

func main() {
	e, err := sh.NewExec(context.Background(), &sh.ExecOptions{
		Storage: &sh.Storage{
			Dir:          "/tmp",
			NotAutoClean: true,
		},
	})
	if err != nil {
		fmt.Printf("new exec fail:%s\n", err)
		return
	}

	dir, _ := os.Getwd()
	fmt.Printf("[%s] %s\n", dir, e.String())

	if err = e.Run("echo hello world"); err != nil {
		fmt.Printf("exec fail:%s\n", err)
	}
}

/*
[go-sh] /bin/bash -ex -o pipefail /tmp/clc9uvco47mm9mrmcbfg
+ echo hello world
hello world
*/
