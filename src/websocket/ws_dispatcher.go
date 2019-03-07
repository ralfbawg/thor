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
	Tasks map[string]*task.WsTask
}

func (m *WsManager) Init() {
	manager = &WsManager{
		Tasks: make(map[string]*task.WsTask,0),
	}
}
func (m *WsManager) CreateTask(appId string) {
	tmp := m.Tasks[appId]
	if tmp == nil {
		map[appId] = &task.WsTask{}
	} else {

	}
}

var (
	upgrader = websocket.Upgrader{
		// 允许跨域
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func WsDispatcher(w http.ResponseWriter, r *http.Request) {
	logging.Debug("ws server start")
	filter.DoFilter(w, r)

	//	var (
	//		wsConn *websocket.Conn
	//		err    error
	//		conn   *connect.Connection
	//		data   []byte
	//		upgrade websocket.Upgrader
	//	)
	//	// 完成ws协议的握手操作
	//	// Upgrade:websocket
	//	if wsConn, err = upgrade.Upgrade(w, r, nil); err != nil {
	//		return
	//	}
	//
	//	if conn, err = connect.InitConnection(wsConn); err != nil {
	//		goto ERR
	//	}
	//ERR:
	//	conn.Close()
}
func verifyAppInfo(r *http.Request) (string, bool) {
	path := strings.Split(r.RequestURI, "/")
	appId := path[2]
	appKey := path[3]
	logging.Debug("app id is %s,app key is %s", appId, appKey)
	//TODO 通过db查询确认
	return appId, true
}
