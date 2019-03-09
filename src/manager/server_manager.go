package manager

import (
	"common/logging"
	"config"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"reflect"
	"strings"
	websocket2 "websocket"
)

const (
	ActionSuffix = "Handler"
)

var upgrade = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// 允许跨域
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
var serverM = &serverManager{}

type serverManager struct {
}

func StartServers() {
	go startHttpServer()
	go startTcpServer()

}
func startHttpServer() {
	logging.Info("start http server")
	tempConfig, _ := config.ConfigStore.GetConfig(false)
	//http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
	//	if conn, err := upgrade.Upgrade(w, r, nil);err!=nil{
	//		fmt.Printf("哦活,error:%s",err)
	//	}else{
	//		//connect.InitConnection(conn)
	//		path := r.RequestURI
	//		websocket2.WsDispatcher(conn,path)
	//	}
	//
	//})
	http.HandleFunc("/", handAdapter)
	//err := http.ListenAndServe("127.0.0.1:8888", nil)
	err := http.ListenAndServe(":"+tempConfig.Server.Port, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
	//http.ListenAndServe(":"+tempConfig.Server.Port, nil)
	logging.Info("http server 启动成功")
}
func startTcpServer() {
	logging.Info("start tcp server")

}
func (c *serverManager) WsHandler(w http.ResponseWriter, r *http.Request) {
	websocket2.WsDispatcher(w,r)


}
func (c *serverManager) ApiHandler(w http.ResponseWriter, r *http.Request) {
	logging.Debug("process api")
}
func handAdapter(w http.ResponseWriter, r *http.Request) {
	logging.Debug("i am in")

	paths := strings.Split(r.RequestURI, "/")
	actionStr := strings.Title(paths[1]) + ActionSuffix
	logging.Debug("action string is " + actionStr)
	obj := reflect.ValueOf(serverM).MethodByName(actionStr)
	if obj.IsValid() {
		obj.Call([]reflect.Value{reflect.ValueOf(w), reflect.ValueOf(r)})
	}
}


