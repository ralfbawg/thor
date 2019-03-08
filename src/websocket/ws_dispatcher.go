package websocket

import (
	"common/logging"
	"github.com/gorilla/websocket"
	"net/http"
	"filter"
	"strings"
	"websocket/task"
)

var manager *WsManager

type WsManager struct {
	tasks          map[string]*task.WsTask
	totalBroadcast chan []byte
}

func (m *WsManager) Init() {
	manager = &WsManager{
		tasks:          make(map[string]*task.WsTask),
		totalBroadcast: make(chan []byte),
	}
	go func() {
		select {
		case msg := <-manager.totalBroadcast:
			for _, v := range manager.tasks {
				v.Broadcast(msg)
			}
		default:

		}
	}()
}

/**
获取当前有多少任务
*/
func (m *WsManager) getSize() int {
	return len(m.tasks)
}

func (m *WsManager) GetOrCreateTask(appId string) *task.WsTask {
	if m.tasks[appId] == nil {
		m.tasks[appId] = task.NewWsTask(appId, m)
	}
	return m.tasks[appId]
}

var upgrade = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// 允许跨域
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func WsDispatcher(w http.ResponseWriter, r *http.Request) {
	logging.Debug("ws server start")
	filter.DoFilter(w, r)
	strings.Split(r.RequestURI, "/")
	if appId, exist := verifyAppInfo(r); exist == true {
		task := manager.GetOrCreateTask(appId)
		if wsConn, err := upgrade.Upgrade(w, r, nil); err != nil {
			return
		} else {
			task.AddClient(appId, wsConn)
			task.Run()
		}

	} else {
		w.Write([]byte("appId 错误或者不存在"))
	}
}

/*
 验证app信息
*/
func verifyAppInfo(r *http.Request) (string, bool) {
	path := strings.Split(r.RequestURI, "/")
	appId := path[2]
	appKey := path[3]
	logging.Debug("app id is %s,app key is %s", appId, appKey)
	//TODO 通过db查询确认
	return appId, appKey == "test"
}
