package main

import (
	"sync"
	"time"
	"github.com/MrMcDuck/xapputil/xlog"
	"github.com/MrMcDuck/xsys/xcron"
)

var wg sync.WaitGroup

func tryFreq(fm *xcron.RateLimiter, id int) {
	defer wg.Add(-1)
	for i := 0; i < 5; i++ {
		fm.MarkAndWaitBlock()
		xlog.Info("Routine %d got a frequency mutex", id)
	}
}

func main() {
	fm, err := xcron.NewRateLimiter(time.Millisecond * 1000)
	if err != nil {
		xlog.Erro(err)
		return
	}

	count := 10
	wg.Add(count)
	for i := 0; i < count; i++ {
		go tryFreq(fm, i)
	}
	wg.Wait()
}
