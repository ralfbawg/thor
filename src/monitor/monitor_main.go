package monitor

import (
	"common"
	"runtime"
	"time"
	"common/logging"
	"runtime/debug"
)

const (
	memCheckInterval = 1 * time.Minute
	byte             = 1
	kb               = 1024 * byte
	mb               = 1024 * kb
	gb               = 1024 * mb
)

type monitorMain struct {
	common.InitI
}

func (m *monitorMain) init() {
	initMemMonitor()
}
func MonitorInit() {
	initMemMonitor()
}

func initMemMonitor() {
	ticker := time.NewTicker(memCheckInterval)
	stats := &runtime.MemStats{}
	go func() {
		for {
			select {
			case <-ticker.C:
				runtime.ReadMemStats(stats)
				inuse := float64(stats.HeapInuse)
				idle := float64(stats.HeapIdle)
				sys := float64(stats.HeapSys)
				logging.Info("当前inuse=%dmb,idle=%dmb,ratio=%d", inuse/mb, idle/mb, idle/sys)
				if (idle/inuse > 2.0 && idle > 500*mb) || idle > 2*gb || (idle/sys > 0.6) {
					logging.Info("应该要返回内存了")
					debug.FreeOSMemory()
				}
			}
		}
	}()

}
