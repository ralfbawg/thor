package main

import (
	"common/logging"
	"net/http"
	"github.com/gorilla"
	"connect"
	"websocket/handler"
)

const ()


func StartServers() {
	logging.Info("start servers")
	startWsServer()

}
func startWsServer() {
	http.HandleFunc("/", wsHandler)
	http.ListenAndServe("0.0.0.0:7777", nil)
	logging.Info("websocket server 启动成功")
}
func startApiServer() {

}


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