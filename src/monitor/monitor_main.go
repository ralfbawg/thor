package monitor

import (
	"common"
	"runtime"
	"time"
	"runtime/debug"
	"common/logging"
)

const (
	memCheckInterval = 1 * time.Minute
	//memCheckInterval = 2 * time.Second
	byte = 1
	kb   = 1024 * byte
	mb   = 1024 * kb
	gb   = 1024 * mb
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
				released := float64(stats.HeapReleased)
				if ((idle-released)/inuse > 2.0 && (idle-released) > 500*mb) || (idle-released) > 2*gb || ((idle-released)/sys > 0.6 && (idle-released) > 500*mb) {
					logging.Info("当前inuse=%fmb,idle=%fmb,sys=%fmb,released=%fmb,ratio=%f", inuse/mb, idle/mb, sys/mb, released/mb, (idle-released)/sys)
					debug.FreeOSMemory()
				}
			}
		}
	}()

}
