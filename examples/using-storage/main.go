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

	e.AddRawCommand([]byte(`

# check user privileges
root_privs() {
  if [[ $(id -u) -eq 0 ]]; then
    return 0
  else
    return 1
  fi
}

# check for required files and deps first
# check if command exists
command_exists() {
  type "${1}" > /dev/null 2>&1
}

if (command_exists ls); then
   echo "command_exists: ls"
fi

`))

	if err = e.Run(); err != nil {
		fmt.Printf("exec fail:%s\n", err)
	}
}

/*
[go-sh] /bin/bash -ex -o pipefail /tmp/clc9uvco47mm9mrmcbfg
+ command_exists ls
+ type ls
+ echo 'command_exists: ls'
command_exists: ls
*/
