package xclock

import (
	"github.com/shirou/gopsutil/host"
	"time"
)

// TODO: support Windows
// NOTICE: It accurates to second only.
func Boottime() (time.Duration, error) {
	uptm, err := host.Uptime()
	if err != nil {
		return 0, err
	}
	return time.Duration(MulDuration(int64(uptm), time.Second)), nil
}
