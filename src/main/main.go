package main

import (
	"common/logging"
	"config"
	"connect"
	"github.com/gorilla"
	"handler"
	"net/http"
)

var (
	upgrader = websocket.Upgrader{
		// 允许跨域
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func wsHandler(w http.ResponseWriter, r *http.Request) {
	//	w.Write([]byte("hello"))
	var (
		wsConn *websocket.Conn
		err    error
		conn   *connect.Connection
		data   []byte
	)
	// 完成ws协议的握手操作
	// Upgrade:websocket
	if wsConn, err = upgrader.Upgrade(w, r, nil); err != nil {
		return
	}

	if conn, err = connect.InitConnection(wsConn); err != nil {
		goto ERR
	}

	handler.AdaptTask()

}

func main() {
	config.Init_main()
	logging.Debug("server start")

}
