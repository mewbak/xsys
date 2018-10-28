package xcron

import (
	"time"
	"github.com/smcduck/xsys/xclock"
)

type SysCron struct {
	triggerTime time.Time
}

func NewSysCron(triggerTime time.Time) (*SysCron, error) {
	return &SysCron{triggerTime: triggerTime}, nil
}

func (sc *SysCron) Wait() {
	for {
		sleepMillis := 100 // default sleep 5000 milliseconds for each loop
		dur := time.Now().Sub(sc.triggerTime)
		durMillis := xclock.NsecToMillis(dur.Nanoseconds())
		if durMillis < 100 {
			break
		}
		time.Sleep(xclock.MillisToDuration(int64(sleepMillis)))
	}
}