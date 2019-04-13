package api

import (
	"encoding/json"
	"game"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"util"
)

const accessKey = "qwerJOQ23j$qw"

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
	clients := game.GameMallInst.Clients()
	list := make([]*ClientInfo, len(clients))
	for _, client := range clients {
		id := client.(*game.GameClient).ID()
		ip := client.(*game.GameClient).IP()
		obj := &ClientInfo{
			ClientId: id,
			ClientIp: ip,
		}
		list = append(list, obj)
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
	clients := game.GameMallInst.Clients()
	resp := &ConnectingUsersResp{
		Code:  0,
		Users: len(clients),
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
	clientId := r.URL.Query().Get("id")
	msg := r.URL.Query().Get("msg")
	myKey := r.URL.Query().Get("key")
	if myKey != accessKey {
		w.Write([]byte("{\"Code\": -1}"))
		return
	}
	obj := game.GameMallInst.Clients()[clientId]
	if obj != nil {
		obj.(*game.GameClient).Send([]byte(msg))
		w.Write([]byte("{\"Code\": 0}"))
	} else {
		w.Write([]byte("{\"Code\": 1}"))
	}
}

func (server *ApiDispatchServer) Broadcast(w http.ResponseWriter, r *http.Request) {
	util.NewConcMap()
}
