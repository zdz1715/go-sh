package main

import (
	"context"
	"fmt"

	"github.com/zdz1715/go-sh"
)

func main() {
	dir := "/var/log"
	e, err := sh.NewExec(context.Background(), &sh.ExecOptions{
		User:    "root",
		WorkDir: dir,
	})
	if err != nil {
		fmt.Printf("new exec fail:%s\n", err)
		return
	}

	fmt.Printf("[%s] %s\n", dir, e.String())

	if err = e.Run("ls | head -n 1"); err != nil {
		fmt.Printf("exec fail:%s\n", err)
	}
}

/*
[/var/log] /bin/bash -ex -o pipefail
+ ls
+ head -n 1
CoreDuet
*/
