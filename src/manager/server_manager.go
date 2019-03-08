package manager

import (
	"common/logging"
	"config"
	"net/http"
	"reflect"
	"strings"
	websocket2 "websocket"
)

const (
	ActionSuffix = "Handler"
)

type serverManager struct {
}

func StartServers() {
	go startHttpServer()
	go startTcpServer()

}
func startHttpServer() {
	logging.Info("start http server")
	tempConfig, _ := config.ConfigStore.GetConfig(false)

	http.HandleFunc("/", handAdapter)
	http.ListenAndServe("127.0.0.1:"+tempConfig.Server.Port, nil)
	logging.Info("http server 启动成功")
}
func startTcpServer() {
	logging.Info("start tcp server")

}
func (c *serverManager) WsHandler(w http.ResponseWriter, r *http.Request) {
	websocket2.WsDispatcher(w, r)
}
func (c *serverManager) ApiHandler(w http.ResponseWriter, r *http.Request) {
	logging.Debug("process api")
}
func handAdapter(w http.ResponseWriter, r *http.Request) {
	logging.Debug("i am in")
	serverM := &serverManager{}
	paths := strings.Split(r.RequestURI, "/")
	actionStr := strings.Title(paths[1]) + ActionSuffix
	logging.Debug("action string is " + actionStr)
	obj := reflect.ValueOf(serverM).MethodByName(actionStr)
	if obj.IsValid() {
		obj.Call([]reflect.Value{reflect.ValueOf(w), reflect.ValueOf(r)})
	}


}
