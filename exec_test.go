package sh

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestExec_Run(t *testing.T) {
	e, err := NewExec(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("exec: %s", e.String())
	_ = e.AddCommand("ls", "-l")
	err = e.Run("echo hello world", "echo golang")
	if err != nil {
		t.Fatal(err)
	}
}

func TestExec_RunWithTimeout(t *testing.T) {
	// 设置超时时间 5s
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	e, err := NewExec(ctx, &ExecOptions{
		//User: "root",
		WorkDir: "/",
		//Storage: &Storage{
		//	Dir:          "/tmp",
		//	NotAutoClean: true,
		//},
		Output: func(num int, line []byte) {
			fmt.Println(string(line))
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	//go func() {
	//	// 也可以手动停止,2s后停止
	//	time.Sleep(2 * time.Second)
	//	e.Cancel()
	//}()

	t.Logf("exec: %s", e.String())
	_ = e.AddCommand("id")
	_ = e.AddCommand(`
docker_ps() {
  docker ps
}
`)
	_ = e.AddCommand("pwd")
	_ = e.AddCommand("cd /usr")
	err = e.Run("docker_ps", "sleep 10&", "echo hello world")
	if err != nil {
		if IsDeadlineExceeded(err) {
			t.Fatal("time out")
		}
		t.Fatal(err)
	}
	t.Logf("last work dir: %s", e.LastWorkDir)
}
