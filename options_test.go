package sh

import (
	"testing"

	"github.com/zdz1715/go-sh/shell"
)

func TestGlobalExecOptionsOverwrite(t *testing.T) {
	t.Logf("%+v\n", GlobalExecOptionsOverwrite(nil))
	t.Logf("%+v\n", GlobalExecOptionsOverwrite(&ExecOptions{
		User:  "root",
		Shell: &shell.Shell{},
	}))
}
