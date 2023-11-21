package sh

import (
	"context"
	"strings"
	"testing"

	"github.com/zdz1715/go-sh/shell"
)

func TestSetDefaultExecOptions(t *testing.T) {
	SetGlobalExecWorkDir("/usr")
	SetGlobalShell(&shell.Shell{
		Type: shell.Sh,
	})
	SetGlobalStorage(&Storage{
		Dir:          "/tmp",
		NotAutoClean: true,
	})

	e1, err := NewExec(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	e1.Run("echo $(pwd)")

	t.Log(strings.Repeat("-", 30))

	e2, err := NewExec(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	e2.Run("echo $(pwd)")

	t.Log(strings.Repeat("-", 30))

	e3, err := NewExec(context.Background(), &ExecOptions{
		WorkDir: "/usr/local",
	})
	if err != nil {
		t.Fatal(err)
	}
	e3.Run("echo $(pwd)")
}
