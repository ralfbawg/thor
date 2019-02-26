package main

import (
	"net/http"
	"time"
	"github.com/gorilla"
	"connect"
	"config"
	"common/logging"
)
func wsHandler(w http.ResponseWriter , r *http.Request){
	//	w.Write([]byte("hello"))
	var(
		wsConn *websocket.Conn
		err error
		conn *connect.Connection
		data []byte
	)
	// 完成ws协议的握手操作
	// Upgrade:websocket
	if wsConn , err = upgrader.Upgrade(w,r,nil); err != nil{
		return
	}

	if conn , err = impl.InitConnection(wsConn); err != nil{
		goto ERR
	}

	// 启动线程，不断发消息
	go func(){
		var (err error)
		for{
			if err = conn.WriteMessage([]byte("heartbeat"));err != nil{
				return
			}
			time.Sleep(1*time.Second)
		}
	}()

	for {
		if data , err = conn.ReadMessage();err != nil{
			goto ERR
		}
		if err = conn.WriteMessage(data);err !=nil{
			goto ERR
		}
	}

ERR:
	conn.Close()

}


func main() {
	config.init.init_main()
	logging.Debug("server start")
	http.HandleFunc("/ws",wsHandler)
	http.ListenAndServe("0.0.0.0:7777",nil)
}
