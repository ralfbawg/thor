package util

import (
	"sync"
	"reflect"
	"time"
)

var (
	BroadcastIns = &BroadcastTask{ //fixme
	}
)

type BroadcastTaskPool struct {
	sync.Pool
	Clean func(i interface{})
}

var BroadcastTaskPoolInst = BroadcastTaskPool{
	Pool: sync.Pool{
		New: func() interface{} {
			return &BroadcastTask{}
		},
	},
	Clean: func(i interface{}) {
		a := reflect.ValueOf(i).MethodByName("Clean")
		if a.IsValid() {
			a.Call([]reflect.Value{})
		}
	},
}

func (p BroadcastTaskPool) Return(item interface{}) {
	p.Clean(item)
	p.Put(item)
}

type BroadcastTask struct {
	appId  string
	taskId int
	uid    string
	msg    []byte
	start  time.Time
}

func (b *BroadcastTask) Clean() {
	b.taskId = 0
	b.appId = ""
	b.uid = ""
	b.start = time.Now()

}
//func (b *BroadcastTask) Run() {
//	m := websocket.GetWsManager()
//	for {
//		m.TaskCount++ //fixme
//	}
//	BroadcastTaskPoolInst.Return(b)
//}

//func NewBroadcastTask(appId string, taskId int, id string, msg []byte) *BroadcastTask {
//	task := BroadcastTaskPoolInst.Get().(*BroadcastTask)
//	ants.Submit(task.Run)
//	return task
//}
