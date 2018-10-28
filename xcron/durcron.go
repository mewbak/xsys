package xcron

import (
	"sync"
	"time"
	"sync/atomic"
	"github.com/pkg/errors"
	"github.com/smcduck/xsys/xclock"
	"github.com/smcduck/xsys/xchan"
)

// Equal duration cron, it's concurrent-safe.
type DurCron struct {
	// User given origin time for cron calc.
	origin time.Time

	// User given duration for cron.
	duration time.Duration

	// Last return trigger time value.
	lastReturn time.Time

	// Mutex for DurCron, avoid repeat return same trigger time in the case of concurrency.
	mu sync.Mutex

	// true: return true in first call Check, false: returns true only duration is reached
	returnTrueFirstCheck bool

	// Flag about whether it is the first time to call Check
	fisrtCheck bool
}

// NewDurCron creates new DurCron object.
// The origin is calc origin time, if it's nil, time.Now() used as origin.
// The d is cron interval.
func NewDurCron(origin *time.Time, returnTrueFirstCheck bool, d time.Duration) *DurCron {
	if origin == nil {
		now := xclock.RoundEarlier(time.Now(), d)
		origin = &now
	}
	c := DurCron{
		origin:               *origin,
		duration:             d,
		lastReturn:           *origin,
		returnTrueFirstCheck: returnTrueFirstCheck,
		fisrtCheck:           true,
	}
	return &c
}

// CheckNowUnblock returns true when now is a trigger point, otherwise false.
// Backlog trigger point history will be ignored, it just return latest trigger point.
func (c *DurCron) CheckNowUnblock() (triggerCount int, trigger bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.fisrtCheck {
		c.fisrtCheck = false
		if c.returnTrueFirstCheck {
			return 1, true
		}
	}

	// Calc float64 times and ceil it to int
	count := int(time.Now().Sub(c.lastReturn) / c.duration)
	//fmt.Println(count, time.Now(), c.lastReturn, time.Now().Sub(c.lastReturn), c.duration)

	if count <= 0 {
		return 0, false
	} else {
		c.lastReturn = c.lastReturn.Add(xclock.MulDuration(int64(count), c.duration))
		return count, true
	}
}





type DurationFunc func(tm time.Time)


type DurationCronWithCallback struct {
	dcron *DurCron
	callback DurationFunc
	runflag atomic.Value
}


func NewDurationCronWithCallback(origin *time.Time,
	startFirstCheck bool, d time.Duration, callback DurationFunc) *DurationCronWithCallback {
		return &DurationCronWithCallback{dcron: NewDurCron(origin, startFirstCheck, d), callback:callback}
}

func (dcc *DurationCronWithCallback) Start() {
	dcc.runflag.Store(1)
	for {
		runflag := dcc.runflag.Load().(int)
		if runflag != 1 {
			time.Sleep(time.Second)
			break
		}
		_, is := dcc.dcron.CheckNowUnblock()
		if is && dcc.callback != nil {
			dcc.callback(time.Now())
		}
		time.Sleep(time.Second)
	}
}

func (dcc *DurationCronWithCallback) Close() error {
	dcc.runflag.Store(0)
	return nil
}




type DurCronCh struct {
	dcron *DurCron
	signal chan struct{}
	done chan struct{}
}

func NewDurCronCh(origin *time.Time, returnTrueFirstCheck bool,
	d time.Duration) *DurCronCh {
	res := DurCronCh{dcron: NewDurCron(origin, returnTrueFirstCheck, d)}

	// TODO: 必须手动指定chan的大小，否则无法正常工作，为啥呢？
	res.signal = make(chan struct{}, 1024)
	res.done = make(chan struct{}, 1)

	go func(dcc *DurCronCh) {
		for {
			select {
			case <-dcc.done:
				break
			default:
				_, is := dcc.dcron.CheckNowUnblock()
				if is {
					dcc.signal <- struct{}{}
				}
				time.Sleep(time.Second)
			}
		}
	}(&res)
	return &res
}

func (dcc *DurCronCh) WaitSignal() error {
	select {
	case <- dcc.signal:
		return nil
	case <- dcc.done:
		return errors.New("Closed cron")
	}
}

func (dcc *DurCronCh) Close() error {
	xchan.SafeCloseChanStruct(dcc.done)
	return nil
}