package shell

import (
	"os/exec"
	"strings"
)

type Type uint8

func (t Type) String() string {
	switch t {
	case Bash:
		return "bash"
	}
	return "sh"
}

const (
	// Bash corresponds to the GNU Bash language, as described in its
	// manual at https://www.gnu.org/software/bash/manual/bash.html.
	//
	// Its string representation is "bash".
	Bash Type = iota

	// Sh
	// https://pubs.opengroup.org/onlinepubs/9699919799/utilities/V3_chap02.html.
	//
	// Its string representation is "posix" or "sh".
	Sh
)

type Shell struct {
	Type  Type
	Set   SetBuiltin
	Unset SetBuiltin
}

func (s Shell) Name() string {
	return s.Type.String()
}

func (s Shell) Path() string {
	name := s.Name()
	if path, err := exec.LookPath(name); err == nil && path != "" {
		return path
	}
	return name
}

func (s Shell) GetFullArgs() []string {
	args := make([]string, 0)
	set := s.Set.Args(false)
	unset := s.Unset.Args(true)
	args = append(args, set...)
	args = append(args, unset...)
	return args
}

func (s Shell) String() string {
	builder := new(strings.Builder)
	builder.WriteString(s.Name())
	builder.Grow(optionMaxLen * 2)
	set := s.Set.String()
	unset := s.Unset.Unset()
	if set != "" {
		builder.WriteByte(' ')
		builder.WriteString(set)
	}
	if unset != "" {
		builder.WriteByte(' ')
		builder.WriteString(unset)
	}
	return builder.String()
}
