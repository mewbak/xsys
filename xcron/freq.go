package xcron

// Frequency limiter designed to limit web api access frequency.
// For example, **.com http api could be visited 3 times in one second,
// you should create a FreqLimiter which duration is (time.second / 3),
// then, MarkAndWaitBlock() before every http request.

import (
	"github.com/anacrolix/sync"
	"time"
	"github.com/MrMcDuck/xapputil/xerror"
)

type FreqLimiter struct {
	d      time.Duration
	last   time.Time
	lastMu sync.Mutex
}

func NewFreqLimiter(d time.Duration) (*FreqLimiter, error) {
	if d < time.Millisecond {
		return nil, xerror.New("Unacceptable too small duration %s for frequency limiter", d.String())
	}
	return &FreqLimiter{d: d}, nil
}

func (f *FreqLimiter) MarkAndWaitUnblock() bool {
	f.lastMu.Lock()
	defer f.lastMu.Unlock()
	if time.Now().Sub(f.last) >= f.d {
		f.last = time.Now()
		return true
	} else {
		return false
	}
}

func (f *FreqLimiter) MarkAndWaitBlock() {
	for {
		time.Sleep(time.Millisecond * 2)
		f.lastMu.Lock()
		if time.Now().Sub(f.last) >= f.d {
			f.last = time.Now()
			f.lastMu.Unlock()
			return
		} else {
			f.lastMu.Unlock()
			continue
		}
	}
}
