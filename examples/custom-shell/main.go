package main

import (
	"fmt"
	"os"

	"github.com/zdz1715/go-sh"
	"github.com/zdz1715/go-sh/shell"
)

func main() {
	e, err := sh.NewExec(&sh.ExecOptions{
		Shell: &shell.Shell{
			Type:  shell.Sh,
			Set:   shell.ErrExit | shell.NoUnset,
			Unset: shell.ErrTrace,
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
[go-sh] /bin/sh -eu +E
hello world
*/
