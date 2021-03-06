package websocket

import (
	"common/logging"
	"github.com/gorilla/websocket"
	"net/http"
	"sync/atomic"
	"task"
	"time"
	"github.com/panjf2000/ants"
	"util"
	"comet/tcp"
	"errors"
	"context"
)

var (
	manager *WsManager
	upgrade = websocket.Upgrader{
		ReadBufferSize:  4096,
		WriteBufferSize: 1024,
		// 允许跨域
		CheckOrigin: func(r *http.Request) bool {
			logging.Debug("i am in checkOrigin")
			//if r.Method != "GET" {
			//	fmt.Println("method is not GET")
			//	return false
			//}
			//if r.URL.Path != "/ws" {
			//	fmt.Println("path error")
			//	return false
			//}
			return true
		},
	}
)

type WsManager struct {
	apps               util.ConcMap
	register           chan *WsApp
	totalBroadcast     chan []byte
	TaskCount          int64
	AppCount           int64
	ClientCount        int64
	CCountC            chan int64
	broadcastTokenPool []int
}
type WsContext struct {
	context.Context
	appId  string
	taskId int
	uid    string
}

func GetWsManager() *WsManager {
	return manager
}

func (m *WsManager) Broadcast(appId string, taskId int, msg []byte) {
	if m.CheckAuth(appId) {
		if appId == "" {
			m.totalBroadcast <- msg
		} else {
			app, ok := m.apps.Get(appId)
			if ok && app != nil {
				app.(*WsApp).Broadcast(msg)
			}
		}
	} else {
		return
	}

}
func (m *WsManager) CheckAuth(appId string) bool {
	return true
}

func WsManagerInit() {
	manager = &WsManager{
		apps:           util.NewConcMap(),
		totalBroadcast: make(chan []byte, 10),
		register:       make(chan *WsApp, 1000),
		CCountC:        make(chan int64, 100),
	}
	tcp.TcpManagerInst.SetWsMethod(WsBroadcast, CloseClient, WsListenersInst.Register, WsListenersInst.Unregister)
	//tcp.TcpManagerInst.SetBroadcast(WsBroadcast)
	//tcp.TcpManagerInst.SetCloseWsHandler(CloseClient)
	//tcp.TcpManagerInst.SetWsListenerRegister(WsListenersInst.Register)
	ants.Submit(func() {
		for {
			select {
			case msg := <-manager.totalBroadcast:
				start := time.Now()
				manager.apps.IterCb(func(key string, v interface{}) {
					v.(*WsApp).Broadcast(msg)
				})
				end := time.Now()
				logging.Debug("broadcast time cost %f second", end.Sub(start).Seconds())
			case app := <-manager.register:
				manager.apps.Set(app.appId, app)
				atomic.AddInt64(&manager.AppCount, 1)
			case ccount := <-manager.CCountC:
				atomic.AddInt64(&manager.ClientCount, ccount)
			}
		}

	})
}

func (m *WsManager) CreateOrGetApp(appId string) (*WsApp, error) {
	if t, exist := m.apps.Get(appId); exist {
		return t.(*WsApp), nil
	} else {
		return NewWsApp(m, appId)
	}
}

/**
获取当前有多少任务
*/
func (m *WsManager) GetSize() int64 {
	return int64(m.apps.Count())
}

/**
获取当前所有任务
*/
func (m *WsManager) GetTasks() int64 {
	count := 0
	for a := range m.apps.IterBuffered() {
		count += len(a.Val.(*WsApp).Tasks)
	}
	return int64(count)
}

/**
获取当前有多少应用
*/
func (m *WsManager) GetApps() map[string]interface{} {
	return m.apps.Items()
}

/**
获取当前有多少任务
*/
func (m *WsManager) GetTaskCount() int64 {
	return atomic.LoadInt64(&m.TaskCount)
}

/**
获取或者新建一个task
*/
func (m *WsManager) GetOrCreateTask(app *WsApp, taskId int) *WsTask {
	if task := app.Tasks[taskId]; task == nil {
		task := NewWsTask(app, taskId)
		app.Tasks[taskId] = task
		return task
	} else {
		return task
	}
}

/**
获取或者新建一个task
*/
func (m *WsManager) GetTask(app *WsApp, taskId int) (*WsTask, error) {
	if task := app.Tasks[taskId]; task != nil {
		return task, nil
	} else {
		return nil, errors.New("没找到task")
	}
}
func (wsManager *WsManager) GetAllClientCount() int64 {
	return atomic.LoadInt64(&wsManager.ClientCount)
}

//手动断开连接
func CloseClient(appId string, taskId int, uid string) {
	if appId == "" {
		return
	} else if app, _ := GetWsManager().apps.Get(appId); app != nil && uid != "" {
		task, err := GetWsManager().GetTask(app.(*WsApp), taskId)
		if err == nil {
			task.GetClient(uid).Close()
		}
	}
}

func WsDispatcher(w http.ResponseWriter, r *http.Request) {
	//logging.Debug("ws server start")
	//filter.DoFilter(w, r)
	param := r.URL.Query()
	//logging.Debug(param.Get("appId"))
	if appId, taskId, uid, exist, error := task.VerifyAppInfo(param); exist == true {
		if conn, err := upgrade.Upgrade(w, r, nil); err != nil {
			logging.Error("哦活,error:%s", err)
		} else {
			if app, err := manager.CreateOrGetApp(appId); err == nil {
				app.AddClient(taskId, uid, conn)
			}

		}
	} else {
		logging.Debug("ws ip(%s)连接错误,签名不过,错误信息(%s)", r.RemoteAddr, error.Error())
		w.Write([]byte("连接错误,签名不过,错误信息(" + error.Error() + ")"))
	}

}

/*
广播主入口
*/
func WsBroadcast(appId string, taskId int, uid string, msg []byte) {
	logging.Debug("ws broadcast param appid(%s),taskId(%d),uid(%s),msg(%s)", appId, taskId, uid, msg)
	if appId == "" {
		return
	} else if app, _ := GetWsManager().apps.Get(appId); app != nil && uid != "" { //单播
		task, err := GetWsManager().GetTask(app.(*WsApp), taskId)
		if err == nil {
			logging.Debug("find ws client uid(%s) exist(%v)", uid, task.GetClient(uid) != nil)
			task.GetClient(uid).Send(msg)
		}
	} else {
		GetWsManager().Broadcast(appId, taskId, msg)
	}
}
