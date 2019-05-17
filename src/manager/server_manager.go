package manager

import (
	"api"
	"common/logging"
	"config"
	"game"
	"log"
	"net/http"
	"reflect"
	"runtime/debug"
	"strings"
	"util"
	websocket2 "websocket"

	"github.com/gorilla/websocket"
	"github.com/panjf2000/ants"
	"net"
	"time"
	"tcp"
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

func (c *serverManager) setWsReadBuffSize(size int) {
	upgrade.ReadBufferSize = size
}
func (c *serverManager) setWsWriteBuffSize(size int) {
	upgrade.WriteBufferSize = size
}
func StartServers() {
	ants.Submit(startHttpServer)
	ants.Submit(startTcpServer)
	ants.Submit(api.ApiDispatchInit)
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
	startShortLinkServer()
	startKeepLiveLinkServer()
	logging.Info("start tcp server")
}
func startShortLinkServer() {

}
func startKeepLiveLinkServer() {
	tempConfig, _ := config.ConfigStore.GetConfig(false)
	tcpAddr, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:"+tempConfig.Server.Tcp.Long)
	tcpListener, _ := net.ListenTCP("tcp", tcpAddr)
	for {
		conn, err := tcpListener.Accept()
		ip := conn.RemoteAddr().String()
		if err != nil {
			continue
		}
		//conn.SetReadDeadline(time.Now().Add(time.Duration(10) * time.Second))
		tcp.MainEntrance(conn,ip)
	}
}
func (c *serverManager) WsHandler(w http.ResponseWriter, r *http.Request) {
	websocket2.WsDispatcher(w, r)

}
func (c *serverManager) ApiHandler(w http.ResponseWriter, r *http.Request) {
	api.ApiDispatch(w, r)

	//logging.Debug("process api")
}

//游戏分发
func (c *serverManager) GameHandler(w http.ResponseWriter, r *http.Request) {
	game.GameDispatch(w, r)
}
func (c *serverManager) DebugHandler(w http.ResponseWriter, r *http.Request) {
	paths := strings.Split(r.RequestURI, "/")
	token := r.URL.Query().Get("token")
	if token != "" || len(paths) < 3 {
		url := "/debug/pprof/"
		http.Redirect(w, r, url, http.StatusFound)
		return
	} else if strings.HasPrefix(paths[2], "mem") {
		util.GetMemoryFile()
	} else if strings.HasPrefix(paths[2], "gc") {
		logging.Info("开始GC")
		debug.FreeOSMemory()
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
