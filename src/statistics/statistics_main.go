package statistics

import (
	"time"
)

var StatArr = []string{"taskSum", "taskClientSum"}
var StatApp []*Statistics

type Statistics struct {
	name   string
	key    string
	subKey string
	count  int
	s      func(a string)
	step   chan int
}

func (s *Statistics) fun(a string) {
	s.count += 1
}
func InitStatistics() {
	t := make([]*Statistics, 100)
	for _, v := range StatArr {
		t = append(t, &Statistics{
			name:  v,
			step:  make(chan int),
			count: 0,
		})
	}
}

func PrintStatistics() {
	tickA := time.NewTicker(20 * time.Second)
	tickB := time.NewTicker(10 * time.Second)
	defer func() {
		tickA.Stop()
		tickB.Stop()
	}()
	for {
		select {
		//case <-tickA.C:
		//	logging.Info("tasks count:%d", GetTaskCount())
		//	//manager := websocket.GetWsManager()
		//	//taskMap := manager.GetTasks()
		//	//taskMap.Foreach(func(s string, i interface{}) {
		//	//	logging.Debug("task key=%s", s)
		//	//	logging.Debug("task value=%s", i.(*websocket.WsTask).GetAppId())
		//	//})
		//case <-tickB.C:
		//	logging.Info("clients count:%d", GetAllClientCount())
		}
	}
}
