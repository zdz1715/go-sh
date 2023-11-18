package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/zdz1715/go-sh"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	e, err := sh.NewExec(ctx)
	if err != nil {
		fmt.Printf("new exec fail:%s\n", err)
		return
	}

	dir, _ := os.Getwd()
	fmt.Printf("[%s] %s\n", dir, e.String())

	e.AddCommand("sleep 10")
	e.AddCommand("echo", "hello world")

	if err = e.Run(); err != nil {
		if sh.IsDeadlineExceeded(err) {
			fmt.Println("error: exec timeout")
		} else {
			fmt.Printf("exec fail:%s\n", err)
		}
	}
}

/*
[go-sh] /bin/bash -ex -o pipefail
+ sleep 10
error: exec timeout
*/
