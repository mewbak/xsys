package xclock

import (
	"time"
	"syscall"
	"github.com/smcduck/xapputil/xerror"
	"github.com/smcduck/xsys/xuser"
)

func SetSystemTimeROOT(t time.Time) error {
		var tv syscall.Timeval
		tv.Sec = t.Unix()
		tv.Usec = 0
		if err := syscall.Settimeofday(&tv); err != nil {
			isAdmin, err2 := xuser.IsRunAsAdmin()
			if err2 == nil && !isAdmin {
				return xerror.New(err.Error() + ", modifying system time requires administrator privileges")
			} else {
				return err
			}
		}
		return nil
}
