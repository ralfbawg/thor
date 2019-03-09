package websocket

import (
	"common/logging"
	"filter"
	"github.com/gorilla/websocket"
	"net/http"
	"strings"
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
		select {
		case msg := <-manager.totalBroadcast:
			for _, v := range manager.tasks {
				v.Broadcast(msg)
			}
		case task := <-manager.register:
			manager.tasks[task.appId] = task
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

//func WsDispatcher(wsConn *websocket.Conn,path string) {
//	logging.Debug("ws server start")
//	filter.DoFilter(w, r)
//	if appId, exist := verifyAppInfo(path); exist == true {
//	task := manager.GetOrCreateTask(appId)
//	task.AddClient(appId, wsConn)
//
//	} else {
//		w.Write([]byte("appId 错误或者不存在"))
//	}
//
//}
func WsDispatcher(w http.ResponseWriter, r *http.Request) {
	logging.Debug("ws server start")
	filter.DoFilter(w, r)
	path := r.RequestURI
	if appId, exist := verifyAppInfo(path); exist == true {
		if conn, err := upgrade.Upgrade(w, r, nil); err != nil {
			logging.Error("哦活,error:%s", err)
		} else {
			task := manager.GetOrCreateTask(appId)
			task.AddClient(appId, conn)
		}

	} else {
		w.Write([]byte("appId 错误或者不存在"))
	}

}

/*
 验证app信息
*/
func verifyAppInfo(path string) (string, bool) {
	paths := strings.Split(path, "/")
	appId := paths[2]
	appKey := paths[3]
	logging.Debug("app id is %s,app key is %s", appId, appKey)
	//TODO 通过db查询确认
	return appId, appKey == "test"
}
