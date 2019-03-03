package main

import (
	"common/logging"
	"config"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	config.Init_main()
	logging.Debug("server start")

	ChanShutdown := make(chan os.Signal)
	signal.Ignore(syscall.SIGHUP)

	<-ChanShutdown
}
