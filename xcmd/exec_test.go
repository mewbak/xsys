package xcmd

import "testing"

func TestExecWait(t *testing.T) {
	ExecWait("ping -c 3 baidu.com", true)
}
