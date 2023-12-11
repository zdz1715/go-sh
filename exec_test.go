package sh

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestExec_Run(t *testing.T) {
	e, err := NewExec()
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
	e, err := NewExecContext(ctx, &ExecOptions{
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
			t.Error("time out")
		}
		t.Error(err)
	}
	t.Logf("finished: %+v lastWorkDir: %s\n", e.finished, e.GetLastWorkDir())
}

func TestExec_Cancel(t *testing.T) {
	e, err := NewExec(&ExecOptions{
		Storage: &Storage{
			Dir: "/tmp",
			//NotAutoClean: true,
		},
		//Shell: &shell.Shell{},
	})
	if err != nil {
		t.Fatal(err)
	}
	go func() {
		// 10s后停止
		time.Sleep(10 * time.Second)
		e.Cancel()
	}()
	// 子进程
	_ = e.AddCommand("sleep 30&")
	err = e.Run("sleep 20")
	if err != nil {
		if IsDeadlineExceeded(err) {
			t.Error("time out")
		}
		t.Error(err)
	}
	t.Logf("finished: %+v lastWorkDir: %s\n", e.finished, e.GetLastWorkDir())
}
