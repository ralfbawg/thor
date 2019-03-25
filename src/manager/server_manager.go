package manager

import (
	"common/logging"
	"config"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"reflect"
	"strings"
	"util"
	websocket2 "websocket"
	"io/ioutil"
	"io"
)

const (
	ActionSuffix = "Handler"
	SuccessMsg   = "success"
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

func (c *serverManager) setWsReadBuffSize(size int) {
	upgrade.ReadBufferSize = size
}
func (c *serverManager) setWsWriteBuffSize(size int) {
	upgrade.WriteBufferSize = size
}
func StartServers() {
	go startHttpServer()
	go startTcpServer()

}
func startHttpServer() {
	logging.Info("start http server")
	tempConfig, _ := config.ConfigStore.GetConfig(false)
	http.HandleFunc("/", handlerAdapter)
	err := http.ListenAndServe(":"+tempConfig.Server.Port, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
	logging.Info("http server 启动成功")
}
func startTcpServer() {
	logging.Info("start tcp server")

}
func (c *serverManager) WsHandler(w http.ResponseWriter, r *http.Request) {
	websocket2.WsDispatcher(w, r)

}
func (c *serverManager) ApiHandler(w http.ResponseWriter, r *http.Request) {
	paths := strings.Split(r.RequestURI, "/")
	if strings.HasPrefix(paths[2], "b") {
		msg := r.URL.Query().Get("msg")
		appId := r.URL.Query().Get("appId")
		uid := r.URL.Query().Get("uid")
		go broadcast(appId, uid, msg)
		w.Write([]byte(SuccessMsg))
	}
	io.Copy(ioutil.Discard, r.Body)
	r.Body.Close()
	//logging.Debug("process api")
}
func (c *serverManager) DebugHandler(w http.ResponseWriter, r *http.Request) {
	paths := strings.Split(r.RequestURI, "/")
	if len(paths) < 3 {
		url := "/debug/pprof/"
		http.Redirect(w, r, url, http.StatusFound)
		return
	} else if strings.HasPrefix(paths[2], "mem") {
		util.GetMemoryFile()
	}

	//logging.Debug("process api")
}
func handlerAdapter(w http.ResponseWriter, r *http.Request) {
	paths := strings.Split(r.RequestURI, "/")
	actionStr := strings.Title(paths[1]) + ActionSuffix

	if strings.Contains(strings.Title(paths[1]), "?") {
		actionStr = strings.Split(strings.Title(paths[1]), "?")[0] + ActionSuffix
	}
	//logging.Debug("action string is " + actionStr)
	obj := reflect.ValueOf(serverM).MethodByName(actionStr)
	if obj.IsValid() {
		obj.Call([]reflect.Value{reflect.ValueOf(w), reflect.ValueOf(r)})
	}
}
func broadcast(appId string, uid string, msg string) {
	logging.Debug("param appid=%s,uid=%s", appId, uid)
	if appId != "" && uid != "" {
		task := websocket2.GetWsManager().GetOrCreateTask(appId)
		task.GetClient(uid).Send([]byte(msg))
	} else {
		websocket2.GetWsManager().Broadcast(appId, msg)
	}
}
