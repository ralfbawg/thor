package game

import (
	"time"
	"common/logging"
	"os"
)

func GameTest() {
	a := make(chan []byte, 2)
	start := time.Now()
	go func() {
		for i := 0; i < 10000; i++ {
			id, _ := CreateOrGetGameRoomId()
			logging.Info("a enter room %d", id)
			//time.Sleep(10 * time.Microsecond)
		}
		a <- []byte("a is finish")
	}()
	go func() {
		for i := 0; i < 10000; i++ {
			id, _ := CreateOrGetGameRoomId()
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
				FindNotEmptyRoom()
				os.Exit(0)
			}
		}
	}
}
