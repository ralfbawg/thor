package websocket

import (
	"common/logging"
	"github.com/gorilla/websocket"
	"net/http"
	"sync/atomic"
	"task"
	"time"
	"util"
	"github.com/panjf2000/ants"
)

var manager *WsManager

var upgrade = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// 允许跨域
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WsManager struct {
	tasks              util.ConcMap
	apps               util.ConcMap
	register           chan *WsApp
	totalBroadcast     chan []byte
	TaskCount          int64
	ClientCount        int64
	broadcastTokenPool []int
}

func GetWsManager() *WsManager {
	return manager
}

func (m *WsManager) Broadcast(appId string, msg string) {
	if m.CheckAuth(appId) {
		return
	}
	if appId == "" {
		m.totalBroadcast <- []byte(msg)
	} else {
		task, ok := m.tasks.Get(appId)
		if ok && task != nil {
			task.(*WsTask).Broadcast([]byte(msg))
		}
	}
}
func (m *WsManager) CheckAuth(appId string) bool {
	return false
}

func WsManagerInit() {
	manager = &WsManager{
		tasks:          util.NewConcMap(),
		totalBroadcast: make(chan []byte, 10),
		register:       make(chan *WsApp, 1000),
	}
	ants.Submit(func() {
		for {
			select {
			case msg := <-manager.totalBroadcast:
				start := time.Now()
				manager.apps.IterCb(func(key string, v interface{}) {
					v.(*WsApp).Broadcast(msg)
				})
				//manager.tasks.Foreach(func(k string, i interface{}) {
				//	i.(*WsTask).Broadcast(msg)
				//})
				end := time.Now()
				logging.Debug("broadcast time cost %f second", end.Sub(start).Seconds())
			case task := <-manager.register:
				manager.apps.Set(task.appId, task)
				atomic.AddInt64(&manager.TaskCount, 1)
			}
		}

	})
}

/**
获取当前有多少任务
*/
func (m *WsManager) GetSize() int64 {
	return int64(m.tasks.Count())
}

/**
获取当前有多少任务
*/
func (m *WsManager) GetTasks() util.ConcMap {
	return m.tasks
}

/**
获取当前有多少任务
*/
func (m *WsManager) GetTaskCount() int64 {
	return atomic.LoadInt64(&m.TaskCount)
}

/**
获取或者新建一个task
*/
func (m *WsManager) GetOrCreateTask(appId string) *WsTask {
	if tmp, ok := m.tasks.Get(appId); !ok || tmp == nil {
		task := NewWsTask(appId, m)
		m.tasks.Set(appId, task)
		return task
	} else {
		return tmp.(*WsTask)
	}
}
func (wsManager *WsManager) GetAllClientCount() int64 {
	return atomic.LoadInt64(&wsManager.ClientCount)
}

func WsDispatcher(w http.ResponseWriter, r *http.Request) {
	//logging.Debug("ws server start")
	//filter.DoFilter(w, r)
	param := r.URL.Query()
	//logging.Debug(param.Get("appId"))
	if appId, id, exist := task.VerifyAppInfo(param); exist == true {
		if conn, err := upgrade.Upgrade(w, r, nil); err != nil {
			logging.Error("哦活,error:%s", err)
		} else {
			task := manager.GetOrCreateTask(appId)
			task.AddClient(id, conn)
		}
	} else {
		w.Write([]byte("appId 错误或者不存在"))
	}

}

/*
广播主入口
*/
func WsBroadcast(appId string, uid string, msg string) {
	logging.Debug("param appid=%s,uid=%s", appId, uid)
	if appId != "" && uid != "" { //单播
		task := GetWsManager().GetOrCreateTask(appId)
		task.GetClient(uid).Send([]byte(msg))
	} else {
		GetWsManager().Broadcast(appId, msg)
	}
}
