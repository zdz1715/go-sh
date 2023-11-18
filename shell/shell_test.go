package shell

import (
	"testing"
)

func TestShell_String(t *testing.T) {
	var sh Shell
	t.Log(sh)
	bashShell := &Shell{
		Type: Bash,
		//Set:   EXPipeFail,
		//Unset: Emacs,
	}
	t.Log(bashShell)
	t.Log(bashShell.GetFullArgs(), len(bashShell.GetFullArgs()))
}

func BenchmarkShell_String(b *testing.B) {
	bashShell := &Shell{
		Type:  Bash,
		Set:   EXPipeFail,
		Unset: Emacs,
	}
	for i := 0; i < b.N; i++ {
		bashShell.String()
	}
}

func BenchmarkShell_GetFullArgs(b *testing.B) {
	bashShell := &Shell{
		Type:  Bash,
		Set:   EXPipeFail,
		Unset: Emacs,
	}
	for i := 0; i < b.N; i++ {
		bashShell.GetFullArgs()
	}
}
