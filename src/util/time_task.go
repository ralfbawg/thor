package util

import "time"

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
