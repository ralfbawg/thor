package main

import (
	"common/logging"
	"config"
	"manager"
	"os"
	"os/signal"
	"syscall"
	"statistics"
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
