package sh

import (
	"fmt"
	"os"

	"github.com/zdz1715/go-sh/shell"
)

var gExecOptions = &ExecOptions{
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
	gExecOptions.WorkDir = dir
}

// SetGlobalExecUser Sets the user for execution globally.
// If the user has been set separately,
// it will not be overwritten.
func SetGlobalExecUser(user string) {
	gExecOptions.User = user
}

// SetGlobalExecOutput Sets the output func for execution globally.
// If the output func has been set separately,
// it will not be overwritten.
func SetGlobalExecOutput(f func(num int, line []byte)) {
	gExecOptions.Output = f
}

// SetGlobalStorage Sets the storage for execution globally.
// If the storage has been set separately,
// it will not be overwritten.
func SetGlobalStorage(storage *Storage) {
	gExecOptions.Storage = storage
}

// SetGlobalIDCreator Sets the ID Creator for execution globally.
// If the ID Creator has been set separately,
// it will not be overwritten.
func SetGlobalIDCreator(creator IDCreator) {
	if creator == nil {
		return
	}
	gExecOptions.IDCreator = creator
}

// SetGlobalShell Sets the shell for execution globally.
// If the shell has been set separately,
// it will not be overwritten.
func SetGlobalShell(shell *shell.Shell) {
	if shell == nil {
		return
	}
	gExecOptions.Shell = shell
}
