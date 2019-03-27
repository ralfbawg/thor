package main

import (
	"common/logging"
	"config"
	"manager"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"statistics"
	"syscall"
)



func main() {
	config.InitMain()
	logging.Debug("server start")
	manager.StartServers()
	statistics.PrintStatistics()
	ChanShutdown := make(chan os.Signal)
	signal.Ignore(syscall.SIGHUP)
	<-ChanShutdown

}
