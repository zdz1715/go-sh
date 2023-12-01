package main

import (
	"fmt"
	"os"

	"github.com/zdz1715/go-sh"
)

func main() {
	e, err := sh.NewExec(&sh.ExecOptions{
		Output: func(num int, line []byte) {
			fmt.Printf("%d| %s\n", num, string(line))
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
[go-sh] /bin/bash -ex -o pipefail
1| + echo hello world
2| hello world
*/
