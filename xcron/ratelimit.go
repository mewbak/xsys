package xcron

// Frequency limiter designed to limit web api access frequency.
// For example, **.com http api could be visited 3 times in one second,
// you should create a RateLimiter which duration is (time.second / 3),
// then, MarkAndWaitBlock() before every http request.

import (
	"github.com/anacrolix/sync"
	"time"
	"github.com/MrMcDuck/xapputil/xerror"
)

type RateLimiter struct {
	d      time.Duration
	last   time.Time
	lastMu sync.Mutex
}

func NewRateLimiter(d time.Duration) (*RateLimiter, error) {
	if d < time.Millisecond {
		return nil, xerror.New("Unacceptable too small duration %s for frequency limiter", d.String())
	}
	return &RateLimiter{d: d}, nil
}

func (r *RateLimiter) MarkAndWaitUnblock() bool {
	r.lastMu.Lock()
	defer r.lastMu.Unlock()
	if time.Now().Sub(r.last) >= r.d {
		r.last = time.Now()
		return true
	} else {
		return false
	}
}

func (r *RateLimiter) MarkAndWaitBlock() {
	for {
		time.Sleep(time.Millisecond * 2)
		r.lastMu.Lock()
		if time.Now().Sub(r.last) >= r.d {
			r.last = time.Now()
			r.lastMu.Unlock()
			return
		} else {
			r.lastMu.Unlock()
			continue
		}
	}
}
