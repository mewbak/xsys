package xcron

import (
	"time"
	"github.com/MrMcDuck/xsys/xclock"
)

type NtpCron struct {
	nc *xclock.NtpClock
	triggerTime time.Time
}

func NewNtpCronONLINE(ntpClock *xclock.NtpClock, triggerTime time.Time) (*NtpCron, error) {
	nc, err := xclock.NewNtpClockONLINE()
	if err != nil {
		return nil, err
	}
	return &NtpCron{nc: nc, triggerTime: triggerTime}, nil
}

func NewNtpCronWithClock(ntpClock *xclock.NtpClock, triggerTime time.Time) (*NtpCron, error) {
	return &NtpCron{nc: ntpClock, triggerTime: triggerTime}, nil
}

func (sc *NtpCron) Wait() {
	for {
		sleepMillis := 100 // default sleep 5000 milliseconds for each loop
		dur := sc.triggerTime.Sub(sc.nc.Now())
		durMillis := xclock.NsecToMillis(dur.Nanoseconds())
		if durMillis < 40 {
			break
		}
		time.Sleep(xclock.MillisToDuration(int64(sleepMillis)))
	}
}