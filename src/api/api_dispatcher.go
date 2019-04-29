package api

import (
	"encoding/json"
	"game"
	"github.com/shirou/gopsutil/cpu"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
	"websocket"
)

var server = new(ApiDispatchServer)
var lastCpuStat = &CpuStat{
	Usage: float64(0),
	Busy:  float64(0),
	Total: float64(0),
}
var cpuStatLock = sync.RWMutex{}

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

func ApiDispatchInit() {
	t := time.NewTicker(3 * time.Second)
	defer t.Stop()
	for {
		if server != nil {
			<-t.C
			cpus, _ := cpu.Percent(3*time.Second, true)
			total := float64(0)
			for _, value := range cpus {
				total += value
			}
			avg := total / float64(len(cpus))
			cpuStatLock.Lock()
			lastCpuStat.Usage = avg
			cpuStatLock.Unlock()
		}
	}
}

//展示在线用户
func (server *ApiDispatchServer) ListOnlineUsers(w http.ResponseWriter, r *http.Request) {
	defer func() {
		io.Copy(ioutil.Discard, r.Body)
		r.Body.Close()
	}()
	tasks := websocket.GetWsManager().GetApps()

	list := make([]*ClientInfo, 0)
	for _, tmp := range tasks {
		app := tmp.(*websocket.WsApp)
		appId := app.GetAppId()
		Tasks := app.Tasks
		for _, task := range Tasks {
			for id, client := range task.GetClients() {
				ip := client.(*websocket.WsTaskClient).GetConn().RemoteAddr().String()
				obj := &ClientInfo{
					AppId:    appId,
					ClientId: id,
					ClientIp: ip,
				}
				list = append(list, obj)
			}
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
	//appId := r.URL.Query().Get("appId")
	//clientId := r.URL.Query().Get("id")
	//msg := r.URL.Query().Get("msg")
	//websocket.WsBroadcast(appId, clientId, msg)
	w.Write([]byte("{\"Code\": 0}"))
}

//广播
func (server *ApiDispatchServer) Broadcast(w http.ResponseWriter, r *http.Request) {
	defer func() {
		io.Copy(ioutil.Discard, r.Body)
		r.Body.Close()
	}()
	paramAppId := r.URL.Query().Get("appId")
	//msg := r.URL.Query().Get("msg")
	if paramAppId != "" {
		//websocket.WsBroadcast(paramAppId, "", msg)
	} else {

		//TODO 新的广播方式
		//tasks := websocket.GetWsManager().GetTasks().Items()
		//for _, tmp := range tasks {
		//	task := tmp.(*websocket.WsTask)
		//	appId := task.GetAppId()
		//	websocket.WsBroadcast(appId, "", msg)
		//}
	}
	w.Write([]byte("{\"Code\": 0}"))
}

// 诊断信息 （内存 CPU）
func (server *ApiDispatchServer) Diagnose(w http.ResponseWriter, r *http.Request) {
	defer func() {
		io.Copy(ioutil.Discard, r.Body)
		r.Body.Close()
	}()
	stats := &runtime.MemStats{}
	runtime.ReadMemStats(stats)
	diagnoseStat := &DiagnoseStat{
		Alloc:    float64(stats.Alloc),
		Inuse:    float64(stats.HeapInuse),
		Idle:     float64(stats.HeapIdle),
		Sys:      float64(stats.HeapSys),
		Released: float64(stats.HeapReleased),
	}
	ret, err := json.Marshal(diagnoseStat)
	if err == nil {
		w.Write(ret)
	} else {
		w.Write([]byte("{}"))
	}
}

func (server *ApiDispatchServer) CpuUsage(w http.ResponseWriter, r *http.Request) {
	defer func() {
		io.Copy(ioutil.Discard, r.Body)
		r.Body.Close()
	}()
	cpuStatLock.RLock()
	cpuStat := &CpuStat{
		Usage: lastCpuStat.Usage,
		Busy:  lastCpuStat.Busy,
		Total: lastCpuStat.Total,
	}
	cpuStatLock.RUnlock()
	ret, err := json.Marshal(cpuStat)
	if err == nil {
		w.Write(ret)
	} else {
		w.Write([]byte("{}"))
	}
}

//游戏
func (server *ApiDispatchServer) Gc(w http.ResponseWriter, r *http.Request) {
	defer func() {
		io.Copy(ioutil.Discard, r.Body)
		r.Body.Close()
	}()
	resulta := ""
	for id, c := range game.GameMallInst.Clients() {
		resulta += "," + id + "|" + c.(*game.GameClient).GetName()
	}
	w.Write([]byte(resulta + "\n"))
	w.Write([]byte(strconv.Itoa(int(game.GameRoomsArr[0]))))
}

//注册websocket应用
func RegisterWs(w http.ResponseWriter, r *http.Request) {

}
