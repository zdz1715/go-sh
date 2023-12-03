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
	"strconv"
	"strings"
	"syscall"

	"github.com/rs/xid"
)

type Exec struct {
	id             string
	xid            string
	lastWorkDir    string // 执行完毕后工作目录位置
	cmd            *exec.Cmd
	ctx            context.Context
	err            error
	file           *os.File
	finished       bool
	finishedRawLen int
	hide           bool
	opts           *ExecOptions

	stdin  io.WriteCloser
	stdout io.ReadCloser
}

func NewExec(execOpts ...*ExecOptions) (*Exec, error) {
	return NewExecContext(context.Background(), execOpts...)
}

func NewExecContext(ctx context.Context, execOpts ...*ExecOptions) (*Exec, error) {
	opts := GlobalExecOptionsOverwrite(execOpts...)

	e := &Exec{
		id:   opts.IDCreator(),
		xid:  xid.New().String(),
		ctx:  ctx,
		opts: opts,
	}

	if e.id == "" {
		return nil, errors.New("id is empty")
	}

	valid, err := CheckStorage(opts.Storage)
	if err != nil {
		return nil, err
	}
	if valid {
		if e.file, err = opts.Storage.CreateFile(e.id); err != nil {
			return nil, err
		}
		e.stdin = e.file
		args := append(opts.Shell.GetFullArgs(), e.file.Name())
		e.cmd = exec.CommandContext(ctx, opts.Shell.Path(), args...)
	} else {
		e.cmd = exec.CommandContext(ctx, opts.Shell.Path(), opts.Shell.GetFullArgs()...)
		if e.stdin, err = e.cmd.StdinPipe(); err != nil {
			return nil, err
		}
	}

	e.cmd.SysProcAttr = &syscall.SysProcAttr{
		// reference：https://jarv.org/posts/command-with-timeout/
		Setpgid: true,
	}

	if opts.User != "" {
		if osUser, err := user.Lookup(opts.User); err != nil {
			return nil, err
		} else {
			uid, _ := strconv.Atoi(osUser.Uid)
			gid, _ := strconv.Atoi(osUser.Gid)
			e.cmd.SysProcAttr.Credential = &syscall.Credential{
				Uid:         uint32(uid),
				Gid:         uint32(gid),
				NoSetGroups: true,
			}
		}
	}

	if opts.WorkDir != "" {
		e.cmd.Dir = opts.WorkDir
	}

	if e.stdout, err = e.cmd.StdoutPipe(); err != nil {
		return nil, err
	}
	// redirect stderr to stdout
	e.cmd.Stderr = e.cmd.Stdout

	return e, nil
}

func (e *Exec) ID() string {
	return e.id
}

func (e *Exec) GetLastWorkDir() string {
	return e.lastWorkDir
}

// Finished returns the value of whether it is finished
func (e *Exec) Finished() bool {
	return e.finished
}

func (e *Exec) String() string {
	if e.cmd != nil {
		return e.cmd.String()
	}
	return ""
}

func (e *Exec) setErr(err error, force bool) {
	if err == nil {
		return
	}
	if force || e.err == nil {
		e.err = &ExecError{
			ID:      e.id,
			Context: e.ctx,
			Err:     err,
		}
	}
}

func (e *Exec) setFinished() {
	if !e.finished {
		e.finished = true
		var err error
		if e.opts.Storage != nil && e.file != nil {
			err = e.opts.Storage.RemoveOrTruncate(e.file, int64(e.finishedRawLen))
			e.setErr(err, false)
		}
		err = e.stdin.Close()
		e.setErr(err, false)
	}
}

// Cancel this execution
func (e *Exec) Cancel() error {
	defer e.setFinished()
	if e.cmd != nil && e.cmd.Process != nil && e.cmd.ProcessState == nil {
		return e.cmd.Process.Kill()
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
	return fmt.Sprintf("%s:%s", e.xid, key)
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
	builder.WriteString("set +x")
	builder.WriteByte('\n')
	builder.WriteString("wait")
	builder.WriteByte('\n')
	builder.WriteString(e.echoKey("pwd:$(pwd)", true))
	builder.WriteByte('\n')
	builder.WriteString(e.echoKey("end"))
	builder.WriteByte('\n')
	raw := builder.Bytes()
	e.finishedRawLen = len(raw)
	return e.AddRawCommand(raw)
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
			e.lastWorkDir = val
			return true
		}
	}

	if e.opts != nil && e.opts.Output != nil &&
		!e.hide && !strings.Contains(line, e.xid) {
		e.opts.Output(num, lineByte)
	}

	return true
}

func (e *Exec) Run(command ...string) error {
	if e.cmd == nil {
		return errors.New("exec: uninitialized")
	}

	if e.finished {
		return errors.New("exec: already finished")
	}

	defer e.setFinished()

	var err error
	if err = e.cmd.Start(); err != nil {
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
		scanner := bufio.NewScanner(e.stdout)
		var num int
		for scanner.Scan() {
			num++
			if !e.parseOutput(num, scanner.Bytes()) {
				break
			}
		}
		if scanner.Err() != nil {
			e.setErr(fmt.Errorf("read: %s", scanner.Err()), false)
		}
	}()

	if err = e.cmd.Wait(); err != nil {
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
