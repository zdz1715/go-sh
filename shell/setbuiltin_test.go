package shell

import (
	"testing"
)

func TestSetBuiltin(t *testing.T) {
	opt := ErrExit | PipeFail | Emacs | XTrace | ErrTrace
	t.Logf("set: %s, unset: %s", opt.String(), opt.Unset())
	t.Logf("set: %s, unset: %s", EXPipeFail.String(), EXPipeFail.Unset())

	only := Emacs | PipeFail
	t.Logf("set: %s, unset: %s", only.String(), only.Unset())
	var empty SetBuiltin
	t.Logf("set: %s, unset: %s", empty.String(), empty.Unset())

}

func BenchmarkSetBuiltin_String(b *testing.B) {
	opt := ErrExit | PipeFail | Emacs | XTrace | ErrTrace
	for i := 0; i < b.N; i++ {
		opt.String()
	}
}

func BenchmarkSetBuiltin_Args(b *testing.B) {
	opt := ErrExit | PipeFail | Emacs | XTrace | ErrTrace
	for i := 0; i < b.N; i++ {
		opt.Args(true)
		//opt.Args(false)
	}
}
