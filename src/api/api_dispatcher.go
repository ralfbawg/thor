package api

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"websocket"
)

var server = new(ApiDispatchServer)

type ApiDispatchServer struct {
}

func ApiDispatch(w http.ResponseWriter, r *http.Request) {
	paths := strings.Split(r.RequestURI, "/")
	actionStr := strings.Title(paths[2])
	if strings.Contains(strings.Title(paths[2]), "?") {
		actionStr = strings.Split(strings.Title(paths[2]), "?")[0]
	}
	obj := reflect.ValueOf(server).MethodByName(actionStr)
	if obj.IsValid() {
		obj.Call([]reflect.Value{reflect.ValueOf(w), reflect.ValueOf(r)})
	}
}

//展示在线用户
func (server *ApiDispatchServer) ListOnlineUsers(w http.ResponseWriter, r *http.Request) {
	defer func() {
		io.Copy(ioutil.Discard, r.Body)
		r.Body.Close()
	}()
	tasks := websocket.GetWsManager().GetTasks().Items()

	list := make([]*ClientInfo, 0)
	for _, tmp := range tasks {
		task := tmp.(*websocket.WsTask)
		appId := task.GetAppId()
		clients := task.GetClients()
		for id, client := range clients {
			ip := client.(*websocket.WsTaskClient).GetConn().RemoteAddr().String()
			obj := &ClientInfo{
				AppId:    appId,
				ClientId: id,
				ClientIp: ip,
			}
			list = append(list, obj)
		}
	}
	resp := &ListOnlineUsersResp{
		Code: 0,
		List: list,
	}
	data, error := json.Marshal(resp)
	if error != nil {
		errorData := string("{\"Code\": 1}")
		w.Write([]byte(errorData))
	} else {
		w.Write([]byte(data))
	}
}

//连接中的用户数
func (server *ApiDispatchServer) ConnectingUsers(w http.ResponseWriter, r *http.Request) {
	defer func() {
		io.Copy(ioutil.Discard, r.Body)
		r.Body.Close()
	}()
	users := websocket.GetWsManager().ClientCount
	resp := &ConnectingUsersResp{
		Code:  0,
		Users: users,
	}
	data, error := json.Marshal(resp)
	if error != nil {
		errorData := string("{\"Code\": 1}")
		w.Write([]byte(errorData))
	} else {
		w.Write([]byte(data))
	}
}

//单播
func (server *ApiDispatchServer) Unicast(w http.ResponseWriter, r *http.Request) {
	defer func() {
		io.Copy(ioutil.Discard, r.Body)
		r.Body.Close()
	}()
	appId := r.URL.Query().Get("appId")
	clientId := r.URL.Query().Get("id")
	msg := r.URL.Query().Get("msg")
	websocket.WsBroadcast(appId, clientId, msg)
	w.Write([]byte("{\"Code\": 0}"))
}

//广播
func (server *ApiDispatchServer) Broadcast(w http.ResponseWriter, r *http.Request) {
	defer func() {
		io.Copy(ioutil.Discard, r.Body)
		r.Body.Close()
	}()
	msg := r.URL.Query().Get("msg")
	tasks := websocket.GetWsManager().GetTasks().Items()
	for _, tmp := range tasks {
		task := tmp.(*websocket.WsTask)
		appId := task.GetAppId()
		websocket.WsBroadcast(appId, "", msg)
	}
	w.Write([]byte("{\"Code\": 0}"))
}
