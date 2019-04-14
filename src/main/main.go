package main

import (
	"common/logging"
	"config"
	"game"
	"github.com/panjf2000/ants"
	"manager"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"statistics"
	"sync"
	"syscall"
	"time"
)

var tmpMap sync.Map

func main() {
	config.InitMain()
	logging.Debug("server start")
	manager.StartServers()
	ants.Submit(statistics.PrintStatistics)
	game.GameMallInst.Init()
	ChanShutdown := make(chan os.Signal)
	signal.Ignore(syscall.SIGHUP)
	<-ChanShutdown

}

func test() {
	a := make(chan []byte, 2)
	start := time.Now()
	go func() {
		for i := 0; i < 10000; i++ {
			id, _ := game.CreateOrGetGameRoomId(0)
			logging.Info("a enter room %d", id)
			//time.Sleep(10 * time.Microsecond)
		}
		a <- []byte("a is finish")
	}()
	go func() {
		for i := 0; i < 10000; i++ {
			id, _ := game.CreateOrGetGameRoomId(0)
			logging.Info("b enter room %d", id)
			//time.Sleep(10 * time.Microsecond)
		}
		a <- []byte("b is finish")
	}()
	count := 0
	for {
		select {
		case msg := <-a:
			logging.Info(string(msg))
			count++
			if count >= 2 {
				end := time.Now()
				logging.Info("total cost time:%s", end.Sub(start).String())
				game.FindNotEmptyRoom()
				os.Exit(0)
			}
		}
	}
}
