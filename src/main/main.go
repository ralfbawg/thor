package main

import (
	"common/logging"
	"config"
	"db"
	"manager"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	config.Init_main()
	logging.Debug("server start")
	db.InitDb()
	manager.StartServers()
	ChanShutdown := make(chan os.Signal)
	signal.Ignore(syscall.SIGHUP)
	<-ChanShutdown
}
