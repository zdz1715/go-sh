package sh

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"

	"github.com/zdz1715/go-sh/shell"
)

type Storage struct {
	Dir          string
	file         *os.File
	NotAutoClean bool
}

func (s *Storage) File() *os.File {
	return s.file
}

type ExecOptions struct {
	IDCreator IDCreator
	Shell     *shell.Shell
	Storage   *Storage
	User      string
	WorkDir   string
	Output    func(num int, line []byte)
}

func defaultExecOptions(execOpts ...*ExecOptions) *ExecOptions {

	opts := new(ExecOptions)
	if len(execOpts) > 0 && execOpts[0] != nil {
		opts = execOpts[0]
	}
	if opts.IDCreator == nil {
		opts.IDCreator = XidCreator
	}
	if opts.Shell == nil {
		opts.Shell = &shell.Shell{
			Type: shell.Bash,
			Set:  shell.EXPipeFail,
		}
	}

	if opts.Output == nil {
		opts.Output = func(num int, line []byte) {
			fmt.Fprintln(os.Stderr, bytesToString(line))
		}
	}

	return opts
}

type Exec struct {
	Context     context.Context
	Cmd         *exec.Cmd
	Options     *ExecOptions
	ID          string
	LastWorkDir string // 执行完毕后工作目录位置

	file         *os.File
	err          error
	hide         bool
	finished     bool
	finishedRaw  []byte
	finishedChan chan struct{}
	stdin        io.WriteCloser
	stdout       io.ReadCloser
}

func NewExec(ctx context.Context, execOpts ...*ExecOptions) (*Exec, error) {
	opts := defaultExecOptions(execOpts...)

	e := &Exec{
		Context:      ctx,
		ID:           opts.IDCreator(),
		Options:      opts,
		finishedChan: make(chan struct{}, 1),
	}
	var err error
	if opts.Storage != nil && opts.Storage.Dir != "" {
		if f, err := os.Stat(opts.Storage.Dir); err != nil {
			return nil, err
		} else if !f.IsDir() {
			return nil, fmt.Errorf("open %s: no such directory", opts.Storage.Dir)
		}

		opts.Storage.file, err = os.Create(filepath.Join(opts.Storage.Dir, e.ID))
		if err != nil {
			return nil, err
		}

		e.stdin = opts.Storage.file
		args := append(opts.Shell.GetFullArgs(), opts.Storage.file.Name())
		e.Cmd = exec.CommandContext(ctx, opts.Shell.Path(), args...)
	} else {
		e.Cmd = exec.CommandContext(ctx, opts.Shell.Path(), opts.Shell.GetFullArgs()...)
	}

	e.Cmd.SysProcAttr = &syscall.SysProcAttr{
		// reference：https://jarv.org/posts/command-with-timeout/
		Setpgid: true,
	}

	if opts.User != "" {
		if osUser, err := user.Lookup(opts.User); err != nil {
			return nil, err
		} else {
			uid, _ := strconv.Atoi(osUser.Uid)
			gid, _ := strconv.Atoi(osUser.Gid)
			e.Cmd.SysProcAttr.Credential = &syscall.Credential{
				Uid:         uint32(uid),
				Gid:         uint32(gid),
				NoSetGroups: true,
			}
		}
	}

	if opts.WorkDir != "" {
		e.Cmd.Dir = opts.WorkDir
	}

	if e.stdin == nil {
		if e.stdin, err = e.Cmd.StdinPipe(); err != nil {
			return nil, err
		}
	}

	if e.stdout, err = e.Cmd.StdoutPipe(); err != nil {
		return nil, err
	}
	// redirect stderr to stdout
	e.Cmd.Stderr = e.Cmd.Stdout

	return e, nil
}

func (e *Exec) String() string {
	if e.Cmd != nil {
		return e.Cmd.String()
	}
	return ""
}

func (e *Exec) setErr(err error, force ...bool) {
	if err == nil {
		return
	}
	if (len(force) > 0 && force[0]) || e.err == nil {
		e.err = &ExecError{
			ID:      e.ID,
			Context: e.Context,
			Err:     err,
		}
	}
}

func (e *Exec) setFinished() {
	if len(e.finishedChan) < cap(e.finishedChan) {
		e.finishedChan <- struct{}{}
		e.finished = true
		// 处理文件
		if e.Options.Storage != nil && e.Options.Storage.File() != nil {
			if !e.Options.Storage.NotAutoClean {
				_ = os.Remove(e.Options.Storage.File().Name())
			} else {
				if stat, err := e.Options.Storage.File().Stat(); err == nil {
					_ = e.Options.Storage.File().Truncate(stat.Size() - int64(len(e.finishedRaw)))
				}
			}
		}
	}
}

// Finished returns the value of whether it is finished
func (e *Exec) Finished() bool {
	return e.finished
}

// Cancel this execution
func (e *Exec) Cancel() error {
	if e.Cmd != nil && e.Cmd.Process != nil && e.Cmd.ProcessState == nil {
		return e.Cmd.Process.Kill()
	}
	return nil
}

func (e *Exec) AddCommand(name string, args ...string) error {
	if name == "" {
		return nil
	}
	n := len(name) + len(args)
	for _, arg := range args {
		n += len(arg)
	}
	builder := new(bytes.Buffer)
	builder.Grow(n)
	builder.WriteString(name)
	for _, arg := range args {
		builder.WriteByte(' ')
		builder.WriteString(arg)
	}
	raw := append(bytes.TrimSpace(builder.Bytes()), '\n')
	return e.AddRawCommand(raw)
}

func (e *Exec) AddRawCommand(raw []byte) error {
	if len(raw) == 0 {
		return nil
	}
	if _, err := e.stdin.Write(raw); err != nil {
		return err
	}
	return nil
}

func (e *Exec) key(key string) string {
	return fmt.Sprintf("%s:%s", e.ID, key)
}

func (e *Exec) echoKey(key string, dbQuote ...bool) string {
	if len(dbQuote) > 0 && dbQuote[0] {
		return fmt.Sprintf(`echo "%s"`, e.key(key))
	}
	return fmt.Sprintf("echo '%s'", e.key(key))
}

func (e *Exec) getKey(key, line string) (val string, found bool) {
	return strings.CutPrefix(line, e.key(key))
}

func (e *Exec) addFinishedRawCommand() error {
	builder := new(bytes.Buffer)
	builder.WriteString(e.echoKey("start"))
	builder.WriteByte('\n')
	builder.WriteString("set + x")
	builder.WriteByte('\n')
	builder.WriteString("wait")
	builder.WriteByte('\n')
	builder.WriteString(e.echoKey("pwd:$(pwd)", true))
	builder.WriteByte('\n')
	builder.WriteString(e.echoKey("end"))
	builder.WriteByte('\n')
	e.finishedRaw = builder.Bytes()
	return e.AddRawCommand(e.finishedRaw)
}

func (e *Exec) parseOutput(num int, lineByte []byte) bool {
	line := bytesToString(lineByte)
	if _, ok := e.getKey("end", line); ok {
		e.setFinished()
		return false
	}

	if _, ok := e.getKey("start", line); ok {
		e.hide = true
		return true
	}

	if e.hide {
		if val, ok := e.getKey("pwd:", line); ok {
			e.LastWorkDir = val
			return true
		}
	}

	if e.Options != nil && e.Options.Output != nil &&
		!e.hide && !strings.Contains(line, e.ID) {
		e.Options.Output(num, lineByte)
	}

	return true
}

func (e *Exec) Run(command ...string) error {
	if e.Cmd == nil {
		return nil
	}
	if e.Finished() {
		return errors.New("exec: already finished")
	}

	var err error
	if err = e.Cmd.Start(); err != nil {
		return err
	}

	for _, s := range command {
		if err = e.AddCommand(s); err != nil {
			return err
		}
	}

	if err = e.addFinishedRawCommand(); err != nil {
		return err
	}

	go func() {
		select {
		case <-e.Context.Done():
			e.setFinished()
		case <-e.finishedChan:
			if closeErr := e.stdin.Close(); closeErr != nil {
				e.setErr(fmt.Errorf("close: %s", closeErr))
			}
		}
	}()

	go func() {
		scanner := bufio.NewScanner(e.stdout)
		var num int
		for scanner.Scan() {
			num++
			if !e.parseOutput(num, scanner.Bytes()) {
				break
			}
		}
		if scanner.Err() != nil {
			e.setErr(fmt.Errorf("read: %s", scanner.Err()))
		}
	}()

	if err = e.Cmd.Wait(); err != nil {
		e.setErr(err, true)
	}

	return e.err
}

type ExecError struct {
	ID      string
	Context context.Context
	Err     error
}

func (e *ExecError) Error() string {
	if e.Context != nil && e.Context.Err() != nil {
		return e.Context.Err().Error()
	}
	return e.Err.Error()
}

func IsDeadlineExceeded(err error) bool {
	return err.Error() == context.DeadlineExceeded.Error()
}
