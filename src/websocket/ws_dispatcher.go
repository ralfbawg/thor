package websocket

import (
	"common/logging"
	"github.com/gorilla/websocket"
	"net/http"
	"net/url"
	"util/uuid"
)

const (
	appIdParam  = "appId"
	appKeyParam = "appKey"
	uidParam    = "id"
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
		m.tasks[appId].Broadcast([]byte(msg))
	}
}

type WsManager struct {
	tasks          map[string]*WsTask
	totalBroadcast chan []byte
	register       chan *WsTask
}

func WsManagerInit() {
	manager = &WsManager{
		tasks:          make(map[string]*WsTask),
		totalBroadcast: make(chan []byte),
		register:       make(chan *WsTask),
	}
	go func() {
		for {
			select {
			case msg := <-manager.totalBroadcast:
				for _, v := range manager.tasks {
					v.Broadcast(msg)
				}
			case task := <-manager.register:
				manager.tasks[task.appId] = task
			}
		}

	}()
}

/**
获取当前有多少任务
*/
func (m *WsManager) getSize() int {
	return len(m.tasks)
}

func (m *WsManager) GetOrCreateTask(appId string) *WsTask {

	if m.tasks[appId] == nil {
		task := NewWsTask(appId, m)
		m.tasks[appId] = task
	}
	return m.tasks[appId]
}

func WsDispatcher(w http.ResponseWriter, r *http.Request) {
	//logging.Debug("ws server start")
	//filter.DoFilter(w, r)
	param := r.URL.Query()
	//logging.Debug(param.Get("appId"))
	if appId,id, exist := verifyAppInfo(param); exist == true {
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
func verifyAppInfo(param url.Values) (string,string, bool) {
	//appId := param.Get(appIdParam)
	//appKey := param.Get(appKeyParam)
	//id := param.Get(uidParam)
	//logging.Debug("app id is %s,app key is %s,uid is %s", appId, appKey, id)
	//TODO 通过db查询确认
	//return id, appKey != "fffasdfasdf" && id != "asdfasdfasd"
	return "test",uuid.Generate().String(), true
}
