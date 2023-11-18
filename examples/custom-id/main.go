package main

import (
	"context"
	"crypto/rand"
	"fmt"

	"github.com/zdz1715/go-sh"
)

func main() {

	e, err := sh.NewExec(context.Background(), &sh.ExecOptions{
		IDCreator: func() string {
			b := make([]byte, 16)
			rand.Read(b)
			return fmt.Sprintf("%x-%x-%x-%x-%x",
				b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
		},
	})

	if err != nil {
		fmt.Printf("new exec fail:%s\n", err)
		return
	}

	fmt.Println("exec id:", e.ID)
}