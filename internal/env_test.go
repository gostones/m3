package internal

import (
	"os"
	"testing"
)

func TestSetDefaultEnviron(t *testing.T) {
	base := "dhnt_base"
	SetDefaultEnviron(base)
	for _, nv := range os.Environ() {
		t.Log(nv)
	}
}
