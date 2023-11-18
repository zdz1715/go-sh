package shell

import (
	"strings"
)

const (
	singleMaxLen = len("-abefhkmnptuvxBCEHPT")
	optionMaxLen = 15 // (Emacs -> Vi)  7 * 2 + 1
)

type SetBuiltin uint32

// Reference: https://www.gnu.org/software/bash/manual/html_node/The-Set-Builtin.html
// set [-abefhkmnptuvxBCEHPT] [-o option-name] [--] [-] [argument …]
const (
	AllExport   SetBuiltin = 1 << iota // -a
	Notify                             // -b
	ErrExit                            // -e
	NoGlob                             // -f
	HashAll                            // -h
	Keyword                            // -k
	Monitor                            // -m
	NoExec                             // -n
	Privileged                         // -p
	OneCmd                             // -t
	NoUnset                            // -u
	Verbose                            // -v
	XTrace                             // -x
	BraceExpand                        // -B
	NoClobber                          // -C
	ErrTrace                           // -E
	HistExpand                         // -H
	Physical                           // -P
	FuncTrace                          // -T

	Emacs     // -o emacs
	History   // -o history
	IgnoreEof // -o ignoreeof
	NoLog     // -o nolog
	PipeFail  // -o pipefail
	Posix     // -o posix
	Vi        // -o vi

)

const (
	EXPipeFail = ErrExit | XTrace | PipeFail
)

func (s SetBuiltin) Args(unset bool) []string {
	symbol := "-"
	if unset {
		symbol = "+"
	}
	o := symbol + "o"
	builder := new(strings.Builder)
	builder.Grow(singleMaxLen)
	builder.WriteString(symbol)
	if s&AllExport != 0 {
		builder.WriteByte('a')
	}
	if s&Notify != 0 {
		builder.WriteByte('b')
	}
	if s&ErrExit != 0 {
		builder.WriteByte('e')
	}
	if s&NoGlob != 0 {
		builder.WriteByte('f')
	}
	if s&HashAll != 0 {
		builder.WriteByte('h')
	}
	if s&Keyword != 0 {
		builder.WriteByte('k')
	}
	if s&Monitor != 0 {
		builder.WriteByte('m')
	}
	if s&NoExec != 0 {
		builder.WriteByte('n')
	}
	if s&Privileged != 0 {
		builder.WriteByte('p')
	}
	if s&OneCmd != 0 {
		builder.WriteByte('t')
	}
	if s&NoUnset != 0 {
		builder.WriteByte('u')
	}
	if s&Verbose != 0 {
		builder.WriteByte('v')
	}
	if s&XTrace != 0 {
		builder.WriteByte('x')
	}
	if s&BraceExpand != 0 {
		builder.WriteByte('B')
	}
	if s&NoClobber != 0 {
		builder.WriteByte('C')
	}
	if s&ErrTrace != 0 {
		builder.WriteByte('E')
	}
	if s&HistExpand != 0 {
		builder.WriteByte('H')
	}
	if s&Physical != 0 {
		builder.WriteByte('P')
	}
	if s&FuncTrace != 0 {
		builder.WriteByte('T')
	}
	// default add symbol, len > 1
	haveSingle := builder.Len() > 1

	args := make([]string, 0, optionMaxLen)

	if haveSingle {
		args = append(args, builder.String())
	}

	if s&Emacs != 0 {
		args = append(args, []string{o, "emacs"}...)
	}
	if s&History != 0 {
		args = append(args, []string{o, "history"}...)
	}
	if s&IgnoreEof != 0 {
		args = append(args, []string{o, "ignoreeof"}...)
	}
	if s&NoLog != 0 {
		args = append(args, []string{o, "nolog"}...)
	}
	if s&PipeFail != 0 {
		args = append(args, []string{o, "pipefail"}...)
	}
	if s&Posix != 0 {
		args = append(args, []string{o, "posix"}...)
	}
	if s&Vi != 0 {
		args = append(args, []string{o, "vi"}...)
	}

	return args
}

func (s SetBuiltin) String() string {
	return strings.Join(s.Args(false), " ")
}

func (s SetBuiltin) Unset() string {
	return strings.Join(s.Args(true), " ")
}
