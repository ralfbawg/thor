package main

import (
	"common/logging"
	"config"
	"manager"
	"os"
	"os/signal"
	"syscall"
)


func main() {
	config.InitMain()
	logging.Debug("server start")
	manager.StartServers()
	ChanShutdown := make(chan os.Signal)
	signal.Ignore(syscall.SIGHUP)
	<-ChanShutdown


}
