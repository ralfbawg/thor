package main

import (
	_ "net/http/pprof"
	"sync"
	"encoding/json"
	"common/logging"
	"os"
	"os/signal"
	"syscall"
	"config"
	"manager"
	"github.com/panjf2000/ants"
	"statistics"
)

var tmpMap sync.Map

func main() {
	config.InitMain()
	logging.Debug("server start")
	manager.StartServers()
	ants.Submit(statistics.PrintStatistics)
	//a, _ := api.GernateAppInfo("onetalk")
	//api.VerifyAppInfo(config.ConfigStore.Ws.App.AppId,"8a6794f01609ffe683fa8a6dbd15af6c7bb76da4",1561365008000,1,"ffc4312164fe038877ad476e453bf74e96a2de8e")
	//game.GameMallInst.Init()
	ChanShutdown := make(chan os.Signal)
	signal.Ignore(syscall.SIGHUP)
	<-ChanShutdown

	//test()
}

func test() {
	//a := "{\"url\":\"abc.com\",\"body\":{\"a\":\"asdasdfsdfsfsfasfd==\"}}"
	//body := []byte(a)
	m := &Message{}
	//m.Body = '{"a":"good"}'
	m.URL = "www.baidu.com"
	//json.Unmarshal(body, m)
	aaa, _ := json.Marshal(m)
	logging.Info("good %s", aaa)

	//a := make(chan []byte, 2)
	//start := time.Now()
	//go func() {
	//	for i := 0; i < 10000; i++ {
	//		id, _ := game.CreateOrGetGameRoomId(0)
	//		logging.Info("a enter room %d", id)
	//		//time.Sleep(10 * time.Microsecond)
	//	}
	//	a <- []byte("a is finish")
	//}()
	//go func() {
	//	for i := 0; i < 10000; i++ {
	//		id, _ := game.CreateOrGetGameRoomId(0)
	//		logging.Info("b enter room %d", id)
	//		//time.Sleep(10 * time.Microsecond)
	//	}
	//	a <- []byte("b is finish")
	//}()
	//count := 0
	//for {
	//	select {
	//	case msg := <-a:
	//		logging.Info(string(msg))
	//		count++
	//		if count >= 2 {
	//			end := time.Now()
	//			logging.Info("total cost time:%s", end.Sub(start).String())
	//			game.FindNotEmptyRoom()
	//			os.Exit(0)
	//		}
	//	}
	//}
}

type Message struct {
	URL  string `json:"url"`
	Body string `json:"body"`
}
