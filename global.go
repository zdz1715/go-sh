package sh

import (
	"fmt"
	"os"

	"github.com/zdz1715/go-sh/shell"
)

var globalExecOptions = &ExecOptions{
	IDCreator: XidCreator,
	Shell: &shell.Shell{
		Type: shell.Bash,
		Set:  shell.EXPipeFail,
	},
	Storage: &Storage{
		Dir:          "",
		NotAutoClean: false,
	},
	User:    "",
	WorkDir: "",
	Output: func(num int, line []byte) {
		fmt.Fprintln(os.Stderr, bytesToString(line))
	},
}

// SetGlobalExecWorkDir Sets the working directory for execution globally.
// If the working directory has been set separately,
// it will not be overwritten.
func SetGlobalExecWorkDir(dir string) {
	globalExecOptions.WorkDir = dir
}

// SetGlobalExecUser Sets the user for execution globally.
// If the user has been set separately,
// it will not be overwritten.
func SetGlobalExecUser(user string) {
	globalExecOptions.User = user
}

// SetGlobalExecOutput Sets the output func for execution globally.
// If the output func has been set separately,
// it will not be overwritten.
func SetGlobalExecOutput(f func(num int, line []byte)) {
	globalExecOptions.Output = f
}

// SetGlobalStorage Sets the storage for execution globally.
// If the storage has been set separately,
// it will not be overwritten.
func SetGlobalStorage(storage *Storage) {
	globalExecOptions.Storage = storage
}

// SetGlobalIDCreator Sets the ID Creator for execution globally.
// If the ID Creator has been set separately,
// it will not be overwritten.
func SetGlobalIDCreator(creator IDCreator) {
	if creator == nil {
		return
	}
	globalExecOptions.IDCreator = creator
}

// SetGlobalShell Sets the shell for execution globally.
// If the shell has been set separately,
// it will not be overwritten.
func SetGlobalShell(shell *shell.Shell) {
	if shell == nil {
		return
	}
	globalExecOptions.Shell = shell
}
