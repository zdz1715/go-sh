package sh

import "github.com/zdz1715/go-sh/shell"

type ExecOptions struct {
	IDCreator IDCreator
	Shell     *shell.Shell
	Storage   *Storage
	User      string
	WorkDir   string
	Output    func(num int, line []byte)
}

func (e *ExecOptions) Copy() *ExecOptions {
	return &ExecOptions{
		IDCreator: e.IDCreator,
		Shell:     e.Shell,
		Storage:   e.Storage,
		User:      e.User,
		WorkDir:   e.WorkDir,
		Output:    e.Output,
	}
}

func GlobalExecOptionsOverwrite(opts ...*ExecOptions) *ExecOptions {
	if len(opts) == 0 || opts[0] == nil {
		return gExecOptions
	}
	eCopy := opts[0].Copy()
	if gExecOptions != nil {
		if eCopy.IDCreator == nil {
			eCopy.IDCreator = gExecOptions.IDCreator
		}
		if eCopy.Shell == nil {
			eCopy.Shell = gExecOptions.Shell
		}
		if eCopy.Storage == nil {
			eCopy.Storage = gExecOptions.Storage
		}
		if eCopy.User == "" {
			eCopy.User = gExecOptions.User
		}
		if eCopy.WorkDir == "" {
			eCopy.WorkDir = gExecOptions.WorkDir
		}
		if eCopy.Output == nil {
			eCopy.Output = gExecOptions.Output
		}
	}
	return eCopy
}
