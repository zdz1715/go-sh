package main

import (
	"context"
	"fmt"
	"os"

	"github.com/zdz1715/go-sh"
)

func main() {
	e, err := sh.NewExec(context.Background())
	if err != nil {
		fmt.Printf("new exec fail:%s\n", err)
		return
	}

	dir, _ := os.Getwd()
	fmt.Printf("[%s] %s\n", dir, e.String())

	e.AddCommand("echo", "pwd: $PWD")
	e.AddCommand("cd", "/")
	// 后台执行
	e.AddCommand("sleep 5 &")

	// 添加多行脚本，配合存储可以生成脚本文件，不调用Run()就不会执行
	e.AddRawCommand([]byte(`
# hello world
hello() {
	echo hello world
}

hello
`))

	if err = e.Run("echo command-1 | cut -d '-' -f2", "echo command2"); err != nil {
		fmt.Printf("exec fail:%s\n", err)
	}

	fmt.Println("exec last work dir:", e.LastWorkDir)
}

/*
[go-sh] /bin/bash -ex -o pipefail
+ cd /
+ hello
+ echo hello world
hello world
+ sleep 5
+ echo command-1
+ cut -d - -f2
1
+ echo command2
command2
exec last work dir: /
*/
