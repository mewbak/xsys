package xproc

import (
	"testing"
)

func TestGetProcInfo(t *testing.T) {
	pid := GetPidOfMyself()
	pi, err := GetProcInfo(pid)
	if err != nil {
		t.Error(err)
		return
	}
}
