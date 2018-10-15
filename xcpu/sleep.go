package xcpu

import (
	"github.com/pkg/errors"
	"time"
)

// 控制CPU使用率，动态调整sleep时间
type DyncSleep struct {
	cpuUsage float64 // 允许的CPU百分比
	lastSleepTime time.Duration
}

func NewDyncSleep(cpuUsage float64) (*DyncSleep, error) {
	if cpuUsage <= 0 || cpuUsage >= 100 {
		return nil, errors.Errorf("Invalid cpuUsage %d", cpuUsage)
	}
	return &DyncSleep{cpuUsage:cpuUsage, lastSleepTime:time.Millisecond}, nil
}

func (s *DyncSleep) Sleep() {
	used, err := GetCombinedUsedPercent(time.Second)
	if err != nil {
		time.Sleep(s.lastSleepTime)
	} else {
		if used > s.cpuUsage {
			s.lastSleepTime += time.Millisecond
		}
		if used < s.cpuUsage {
			s.lastSleepTime -= time.Millisecond

		}
		if s.lastSleepTime <= 0 {
			s.lastSleepTime = time.Millisecond
		}
		time.Sleep(s.lastSleepTime)
	}
}
