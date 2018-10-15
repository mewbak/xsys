package xclock

// A clock who is independent with system clock and sync with NTP server.
// Change system clock need ROOT, but NtpClock doesn't.

import (
	"sync"
	"time"
)

// github.com/hlandau/degoutils/clock

type NtpClock struct {
	diff     time.Duration
	diffRwmu sync.RWMutex
}

func NewNtpClockONLINE() (*NtpClock, error) {
	ntptime, err := GetNetTimeInLocalONLINE()
	if err != nil {
		return nil, err
	}

	nc := NtpClock{
		diff: ntptime.Sub(time.Now()),
	}
	go nc.syncroutine()
	return &nc, nil
}

func (nc *NtpClock) syncroutine() {
	for {
		time.Sleep(time.Second * 5)
		ntptime, err := GetNetTimeInLocalONLINE()
		if err != nil {
			continue
		}
		nc.diffRwmu.Lock()
		nc.diff = ntptime.Sub(time.Now())
		nc.diffRwmu.Unlock()
	}
}

func (nc *NtpClock) Now() time.Time {
	nc.diffRwmu.RLock()
	diff := nc.diff
	nc.diffRwmu.RUnlock()
	return time.Now().Add(diff)
}
