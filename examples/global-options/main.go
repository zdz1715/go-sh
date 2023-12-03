package main

import (
	"fmt"

	"github.com/zdz1715/go-sh/shell"

	"github.com/zdz1715/go-sh"
)

func main() {
	// 设置全局执行的工作目录
	sh.SetGlobalExecWorkDir("/usr")
	// 设置全局执行的用户
	sh.SetGlobalExecUser("")
	// 设置全局执行的输出方法
	sh.SetGlobalExecOutput(func(num int, line []byte) {
		fmt.Println(string(line))
	})
	// 设置全局执行脚本存储方式
	sh.SetGlobalStorage(&sh.Storage{
		Dir: "/tmp",
		//NotAutoClean: true,
	})
	// 设置全局的id生成方式
	sh.SetGlobalIDCreator(func() string {
		return "jenkins-" + sh.XidCreator()
	})

	// 设置全局的shell 类型
	sh.SetGlobalShell(&shell.Shell{
		Type: shell.Sh,
	})

	e1, _ := sh.NewExec()
	e2, _ := sh.NewExec()
	e3, _ := sh.NewExec()
	e1.Run("echo e1: $(pwd)")
	e2.Run("echo e2: $(pwd)")
	e3.Run("echo e3: $(pwd)")
}
