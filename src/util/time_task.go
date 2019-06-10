package util

import (
	"time"
	"github.com/panjf2000/ants"
)

var (
	wrPoolExtendFactor = 0.8
	wrPoolDefaultSize  = 10000
	wrPool, _          = ants.NewPool(wrPoolDefaultSize)
	funcs              = make([]func(), 0)
)

func TimeTask(t time.Duration, count int, f func(n int)) {
	go func() {
		tick := time.NewTicker(t)
		initCount := 0
		for range tick.C {
			if initCount < count {
				f(initCount)
				initCount++
			} else {
				tick.Stop()
			}
		}
	}()
}
func SubmitTaskAndResize(pool *ants.Pool, defaultSize int, factor float64, f ...func()) {
	if float64(pool.Running())/float64(defaultSize) > factor {
		defaultSize *= 2
		pool.Tune(defaultSize)
	}
	for _, v := range f {
		pool.Submit(v)
	}
}
