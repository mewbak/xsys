package xclock

import (
	"time"
	"syscall"
	"github.com/smcduck/xsys/xuser"
	"github.com/pkg/errors"
)

func SetSystemTimeROOT(t time.Time) error {
		var tv syscall.Timeval
		tv.Sec = t.Unix()
		tv.Usec = 0
		if err := syscall.Settimeofday(&tv); err != nil {
			isAdmin, err2 := xuser.IsRunAsAdmin()
			if err2 == nil && !isAdmin {
				return errors.Wrap(err, "modifying system time requires administrator privileges")
			} else {
				return err
			}
		}
		return nil
}
