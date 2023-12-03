package sh

import (
	"testing"
)

func TestCheckStorage(t *testing.T) {
	t.Log(CheckStorage(nil))
	t.Log(CheckStorage(&Storage{
		Dir: "",
	}))
	t.Log(CheckStorage(&Storage{
		Dir: "/tmp",
	}))
	t.Log(CheckStorage(&Storage{
		Dir: "/tmp/111",
	}))
}
