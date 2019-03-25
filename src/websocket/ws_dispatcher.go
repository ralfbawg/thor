package websocket

import (
	"github.com/gorilla/websocket"
	"net/http"
	"net/url"
	"util/uuid"
	"util"
	"common/logging"
	"sync/atomic"
	"time"
)

const (
	appIdParam  = "appId"
	appKeyParam = "appKey"
	uidParam    = "uid"
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

func GetWsManager() *WsManager {
	return manager
}

func (m *WsManager) Broadcast(appId string, msg string) {
	if appId == "" {
		m.totalBroadcast <- []byte(msg)
	} else {
		task := m.tasks.Get(appId)
		if task != nil {
			task.(*WsTask).Broadcast([]byte(msg))
		}
	}
}

type WsManager struct {
	tasks          *util.ConcurrentMap
	totalBroadcast chan []byte
	register       chan *WsTask
	TaskCount      int64
	ClientCount    int64
}

func WsManagerInit() {
	manager = &WsManager{
		tasks:          util.NewConcurrentMap(),
		totalBroadcast: make(chan []byte),
		register:       make(chan *WsTask, 1000),
	}
	go func() {
		for {
			select {
			case msg := <-manager.totalBroadcast:
				start := time.Now()
				manager.tasks.Foreach(func(k string, i interface{}) {
					i.(*WsTask).Broadcast(msg)
				})
				end := time.Now()
				logging.Debug("broadcast time cost %f second", end.Sub(start).Seconds())
			case task := <-manager.register:
				manager.tasks.Put(task.appId, task)
				atomic.AddInt64(&manager.TaskCount, 1)
			}
		}

	}()
}

/**
获取当前有多少任务
*/
func (m *WsManager) GetSize() int64 {
	return int64(len(m.tasks.Map))
}

/**
获取当前有多少任务
*/
func (m *WsManager) GetTasks() *util.ConcurrentMap {
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

	if m.tasks.Get(appId) == nil {
		task := NewWsTask(appId, m)
		m.tasks.Put(appId, task)
	}
	return m.tasks.Get(appId).(*WsTask)
}
func (wsManager *WsManager) GetAllClientCount() int64 {
	return atomic.LoadInt64(&wsManager.ClientCount)
}

func WsDispatcher(w http.ResponseWriter, r *http.Request) {
	//logging.Debug("ws server start")
	//filter.DoFilter(w, r)
	param := r.URL.Query()
	//logging.Debug(param.Get("appId"))
	if appId, id, exist := verifyAppInfo(param); exist == true {
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
 验证app信息
*/
func verifyAppInfo(param url.Values) (string, string, bool) {
	appId := param.Get(appIdParam)
	//appKey := param.Get(appKeyParam)
	id := param.Get(uidParam)
	if id == "" {
		id = uuid.Generate().String()
	}
	//logging.Debug("app id is %s,app key is %s,uid is %s", appId, appKey, id)
	//TODO 通过db查询确认
	//return id, appKey != "fffasdfasdf" && id != "asdfasdfasd"
	return appId, id, true
}
