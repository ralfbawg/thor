package api

import (
	"encoding/json"
	"game"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
)

const SuccessMsg = "success"

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
		errorData := string("{\"code\": 1}")
		w.Write([]byte(errorData))
	} else {
		w.Write([]byte(data))
	}
}
