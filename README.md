# go-sh
go写的shell脚本处理包，为提高工作效率而开发
- 支持执行多条命令，并等待所有后台任务结束
- 自定义实时输出方式
- 自定义执行ID生成方式，便于追踪执行记录
- 可快捷指定shell类型和[Set-Builtin](https://www.gnu.org/software/bash/manual/html_node/The-Set-Builtin.html)
- 支持获取执行完所在的工作目录，便于设置下一次执行的工作目录
- 支持根据命令生成脚本文件去执行，可存储每次执行脚本

## Contents
- [Installation](#Installation)
- [Quick start](#quick-start)
- [Examples](#examples)
## Installation
```shell
go get -u github.com/zdz1715/go-sh@latest
```

## Quick start
```go
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

```
## Examples
- [Exec with timeout](./examples/timeout/main.go)
- [Using user, dir](./examples/using-user-dir/main.go)
- [Using Storage](./examples/using-storage/main.go)
- [Custom Shell](./examples/custom-shell/main.go)
- [Custom Output](./examples/custom-output/main.go)
- [Custom ID](./examples/custom-id/main.go)

