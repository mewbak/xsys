package xfs

import (
	"testing"
)

func TestMakeDir(t *testing.T) {
	if err := MakeDir("abc"); err != nil {
		t.Error(err)
		return
	}
	if err := RemoveDir("abc"); err != nil {
		t.Error(err)
		return
	}
}
